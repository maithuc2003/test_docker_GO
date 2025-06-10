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

    stage('Merge Develop into Main') {
      steps {
        script {
          bat 'git config user.name "Jenkins"'
          bat 'git config user.email "jenkins@example.com"'
          bat 'git fetch origin develop'
          bat 'git merge origin/develop'
          
          // Nếu muốn cập nhật lên GitHub:
          withCredentials([usernamePassword(credentialsId: 'GITHUB_username_with_password', usernameVariable: 'GIT_USER', passwordVariable: 'GIT_PASS')]) {
            bat '''
              git push https://%GIT_USER%:%GIT_PASS%@github.com/maithuc2003/test_docker_GO.git main
            '''
          }
        }
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
