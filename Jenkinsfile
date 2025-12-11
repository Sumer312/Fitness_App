pipline {
  agent any
  tools {
  }
  stages("build"){
    echo "Generating templ"
    sh "templ generate"
    echo "Compling project"
    sh "go build -o Fitness"
    echo "Runing the executable"
    sh "./Fitness"
  }
}
