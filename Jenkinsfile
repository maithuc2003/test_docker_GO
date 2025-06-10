pipeline {
  agent any

  environment {
    DOCKER_IMAGE = 'maithuc2003/go-book-api'
    DOCKER_TAG = 'latest'
  }

  stages {
    stage('Build Docker Image') {
      steps {
        bat "docker build -t %DOCKER_IMAGE%:%DOCKER_TAG% ."
      }
    }

    stage('Push Docker Image') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'DOCKERHUB', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
          bat """
          echo %PASS% | docker login -u %USER% --password-stdin
          docker push %DOCKER_IMAGE%:%DOCKER_TAG%
          """
        }
      }
    }
  }
}
