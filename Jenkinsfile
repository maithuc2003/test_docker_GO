pipeline {
  agent any

  environment {
    DOCKER_IMAGE = 'maithuc2003/go-book-api'
    DOCKER_TAG = "latest"
  }

  stages {
    stage('Checkout') {
      steps {
        git url: 'https://github.com/maithuc2003/test_docker_GO.git', branch: 'main'
      }
    }
    stage('Merge Feature Branch') {
      steps {
        script {
          bat 'git fetch origin feature-branch'
          bat 'git merge origin/feature-branch --no-ff -m "Merge feature-branch into main"'
        }
      }
    }
    stage('Build Docker Image') {
      steps {
        script {
          bat "docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} ."
        }
      }
    }
    stage('Push Docker Image') {
      steps {
        script {
          withCredentials([usernamePassword(credentialsId: 'dockerhub-creds', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
            bat "docker login -u %USER% -p %PASS%"
            bat "docker push ${DOCKER_IMAGE}:${DOCKER_TAG}"
          }
        }
      }
    }
  } // <-- đóng stages
} // <-- đóng pipeline
