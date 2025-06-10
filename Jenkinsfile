pipeline {
    agent any

    environment {
        DOCKER_IMAGE = 'maithuc2003/go-book-api'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Docker Image') {
            steps {
                bat "docker build -t %DOCKER_IMAGE%:%BRANCH_NAME% ."
                bat "docker tag %DOCKER_IMAGE%:%BRANCH_NAME% %DOCKER_IMAGE%:latest"
            }
        }

        stage('Push Docker Image') {
            when {
                branch 'main'
            }
            steps {
                withCredentials([usernamePassword(credentialsId: 'DOCKERHUB', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    bat """
                        docker login -u %DOCKER_USER% -p %DOCKER_PASS%
                        docker push %DOCKER_IMAGE%:latest
                    """
                }
            }
        }

        stage('Deploy with Docker Compose') {
            when {
                branch 'main'
            }
            steps {
                withCredentials([file(credentialsId: 'my-env-file_new', variable: 'ENV_PATH')]) {
                    bat """
                        copy %ENV_PATH% .env
                         docker-compose -f docker-compose.yaml down
                        docker-compose -f docker-compose.yaml up -d --build
                    """
                }
            }
        }
    }

    post {
        always {
            bat 'del .env' // Xoá .env sau khi chạy
        }
    }
}
