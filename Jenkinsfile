
import altipla.CI


node {
  def ci = new CI()

  ci.init this

  stage('Checkout') {
    checkout scm
  }

  stage('king') {
    dir('tools') {
      sh 'actools go build -o king ./cmd/king'
      ci.gsutil "-h 'Cache-Control: no-cache' cp king gs://tools.altipla.consulting/bin/king"
    }
  }
}
