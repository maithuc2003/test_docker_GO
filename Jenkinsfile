pipeline {
  agent any

  environment {
    DOCKER_IMAGE = 'yourdockerhubusername/your-image-name'
    DOCKER_TAG = "latest"
  }

  stages {
    stage('Checkout') {
      steps {
        // Checkout branch target (ví dụ main)
        git url: 'https://github.com/maithuc2003/test_docker_GO.git', branch: 'main'
      }
    }

    stage('Merge Feature Branch') {
      steps {
        script {
          // Merge branch feature-branch vào main local
          sh 'git fetch origin feature-branch'
          sh 'git merge origin/feature-branch --no-ff -m "Merge feature-branch into main"'
        }
      }
    }

    stage('Build Docker Image') {
      steps {
        script {
          sh "docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} ."
        }
      }
    }

    stage('Push Docker Image') {
      steps {
        script {
          withCredentials([usernamePassword(credentialsId: 'dockerhub-creds', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
            sh "echo $PASS | docker login -u $USER --password-stdin"
            sh "docker push ${DOCKER_IMAGE}:${DOCKER_TAG}"
          }
        }
      }
    }
  }
}
