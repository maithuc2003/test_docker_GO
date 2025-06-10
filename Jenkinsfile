pipeline {
  agent any

  environment {
    DOCKER_IMAGE = 'maithuc2003/go-book-api'
    DOCKER_TAG = "latest"
  }

  stages {
    stage('Checkout') {
      steps {
        git url: 'https://github.com/maithuc2003/test_docker_GO.git', branch: 'main', credentialsId: 'GITHUB_username_with_password'
      }
    }

    stage('Build Docker Image') {
      steps {
        bat "docker build -t %DOCKER_IMAGE%:%DOCKER_TAG% ."
      }
    }

    stage('Push Docker Image') {
      steps {
        withCredentials([usernamePassword(credentialsId: '855e4714-4ec0-426c-b7ef-faaec4463a7d', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
          bat '''
            docker login -u %USER% -p %PASS%
            docker push %DOCKER_IMAGE%:%DOCKER_TAG%
          '''
        }
      }
    }
  }
}
