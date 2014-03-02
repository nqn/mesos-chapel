package main

import (
    "io/ioutil"
    "fmt"
    "strconv"
    "flag"
    "log"
    "os"
    "io"
    "os/exec"
    "strings"
    "mesos.apache.org/mesos"
    "code.google.com/p/goprotobuf/proto"
)

func triggerAgents(chapel_program []string, agents map[string]uint64) {
  f, err := ioutil.TempFile("/tmp", "chapel")
  if err != nil {
    log.Fatal(err)
  }

  for hostname, port := range agents {
    line := hostname + ":" + strconv.FormatUint(port, 10) + "\n"
    f.Write([]byte(line))
  }

  f.Close()

  pwd, _ := os.Getwd()

  name, _ := ioutil.TempDir("/tmp", "chapel")
  os.Chdir(name)

  // TODO(nnielsen): Assumes program + program_real executables.
  os.Symlink(pwd + "/" + chapel_program[0], name + "/" + chapel_program[0])
  os.Symlink(pwd + "/" + chapel_program[0] + "_real", name + "/" + chapel_program[0] + "_real")
  err = os.Symlink(pwd + "/bin/chapel-client", name + "/chapel-client")
  if err != nil {
    log.Fatal(err.Error())
  }

  cmd := exec.Command(chapel_program[0], chapel_program[1:]...)

  stdout, _ := cmd.StdoutPipe()
  go io.Copy(os.Stdout, stdout)

  stderr, _ := cmd.StderrPipe()
  go io.Copy(os.Stderr, stderr)

  cmd.Env = os.Environ()
  cmd.Env = append(cmd.Env, "SSH_CMD=./chapel-client")
  cmd.Env = append(cmd.Env, "SSH_OPTIONS=" + f.Name())

  if err := cmd.Run(); err != nil {
    log.Fatal(err)
  }

  cmd.Wait()
}

func main() {
  finishedLocales := 0
  runningLocales := 0
  taskId := 0
  exit := make(chan bool)
  agents := make(map[string]uint64)
  saturated := false

  pwd, _ := os.Getwd()

  master := flag.String("master", "localhost:5050", "Location of leading Mesos master")
  bootstrap := flag.String("bootstrap", pwd + "/chapel-bootstrap.tgz", "Location of bootstrap package")
  locales := flag.Int("locales", 1, "Number of Chapel locales i.e. number of nodes")
  flag.Parse()
  chapel_program := flag.Args()

  args := len(chapel_program)
  if args == 0 {
    fmt.Println("No Chapel program found")
    fmt.Println("Syntax: ./chapel (<options>) <chapel-program>")
    os.Exit(1)
  }

  chapel_program = append(chapel_program, "-nl " + strconv.Itoa(*locales))

  driver := mesos.SchedulerDriver {
    Master: *master,
    Framework: mesos.FrameworkInfo {
        Name: proto.String("Chapel: " + strings.Join(chapel_program, " ")),
        User: proto.String(""),
    },

    Scheduler: &mesos.Scheduler {
      ResourceOffers: func(driver *mesos.SchedulerDriver, offers []mesos.Offer) {
        for _, offer := range offers {
          if saturated {
            driver.DeclineOffer(offer.Id)
            continue
          }

          var port uint64 = 0
          var cpus float64 = 0
          var mem float64 = 0

          for _, resource := range offer.Resources {
            if resource.GetName() == "cpus" {
              cpus = *resource.GetScalar().Value
            }

            if resource.GetName() == "mem" {
              mem = *resource.GetScalar().Value
            }

            if resource.GetName() == "ports" {
              r := (*resource).GetRanges().Range[0]
              port = r.GetBegin()
            }
          }

          agents[*offer.Hostname] = port

          command := &mesos.CommandInfo {
            Value: proto.String("./bin/chapel-agent -port=" + strconv.FormatUint(port, 10)),
            Uris:  []*mesos.CommandInfo_URI {
              &mesos.CommandInfo_URI { Value: bootstrap },
            },
          }

          taskId++
          tasks := []mesos.TaskInfo {
            mesos.TaskInfo {
              Name:    proto.String("Chapel-agent"),
              TaskId:  &mesos.TaskID {
                Value: proto.String("Chapel-agent-" + strconv.Itoa(taskId)),
              },
              SlaveId: offer.SlaveId,
              Command: command,
              Resources: []*mesos.Resource {
                mesos.ScalarResource("cpus", cpus),
                mesos.ScalarResource("mem", mem),
                mesos.RangeResource("ports", port, port),
              },
            },
          }

          driver.LaunchTasks(offer.Id, tasks)
        }
      },

      StatusUpdate: func(driver *mesos.SchedulerDriver, status mesos.TaskStatus) {
        if (*status.State == mesos.TaskState_TASK_FINISHED) {
          finishedLocales++
          if finishedLocales >= *locales {
            exit <- true
          }
        } else if (*status.State == mesos.TaskState_TASK_RUNNING) {
          runningLocales++
          if runningLocales >= *locales {
            saturated = true
            go triggerAgents(chapel_program, agents)
          } else {
            fmt.Println("[" , runningLocales , "/" , *locales , "] Setting up locale..")
          }
        } else {
          fmt.Println("Received task status: " + mesos.TaskState_name[int32(*status.State)])
          if (status.Message != nil) {
            fmt.Println("Message: " + *status.Message)
          }
        }
      },
    },
  }

  driver.Init()
  defer driver.Destroy()

  driver.Start()
  <-exit
  driver.Stop(false)
}
