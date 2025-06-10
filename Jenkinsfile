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
                bat "docker build -t ${DOCKER_IMAGE}:${env.BRANCH_NAME} ."
            }
        }

        stage('Push Docker Image') {
            when {
                branch 'main'   // chỉ chạy stage này nếu đang ở branch main
            }
            steps {
                withCredentials([string(credentialsId: 'DOCKERHUB', variable: 'PASS')]) {
                    bat """
                    echo %PASS% | docker login -u maithuc2003 --password-stdin
                    docker push ${DOCKER_IMAGE}:${env.BRANCH_NAME}
                    """
                }
            }
        }
    }
}
