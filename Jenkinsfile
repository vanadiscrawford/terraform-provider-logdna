library 'magic-butler-catalogue'
def PROJECT_NAME = 'terraform-provider-logdna'
def CURRENT_BRANCH = currentBranch()
def MAIN_BRANCH = 'main'
def GIT_REPO = 'logdna/terraform-provider-logdna'
def GIT_AUTHOR = 'logdnabot'
def TRIGGER_PATTERN = ".*@logdnabot.*"

pipeline {
  agent none

  options {
    timestamps()
    ansiColor 'xterm'
  }

  triggers {
    issueCommentTrigger(TRIGGER_PATTERN)
  }

  stages {
    stage('Validate PR Source') {
      when {
        expression { env.CHANGE_FORK }
        not {
          triggeredBy 'issueCommentCause'
        }
      }

      steps {
        error('A maintainer needs to approve this PR for CI by commenting')
      }
    }

    stage('Test') {
      agent {
        node {
          label 'ec2-fleet'
          customWorkspace "${PROJECT_NAME}-${BUILD_NUMBER}"
        }
      }

      environment {
        GIT_BRANCH = "${CURRENT_BRANCH}"
      }

      steps {
        script {
          withCredentials([
            string(credentialsId: 'logdna-gpg-key', variable: 'GPG_KEY'),
            string(credentialsId: 'terraform-provider-servicekey', variable: 'SERVICE_KEY'),
            string(credentialsId: 'terraform-provider-coveralls', variable: 'COVERALLS_TOKEN'),
            string(credentialsId: 'terraform-test-s3-bucket', variable: 'S3_BUCKET'),
            string(credentialsId: 'terraform-test-gcs-bucket', variable: 'GCS_BUCKET'),
            string(credentialsId: 'terraform-test-gcs-projectid', variable: 'GCS_PROJECTID')
          ]) {
            sh '''
              set +x
              echo "$GPG_KEY" > gpgkey.asc
              make postcov
              make test-release
            '''
          }

          sh '''
            set +x
            git checkout -b ${GIT_BRANCH} origin/${GIT_BRANCH}
            git fetch --tags
            export CURRENT_TAG=$(make version-current)
            export NEXT_TAG=$(make version-next)
            echo "Latest: ${CURRENT_TAG}"
            echo "Next: ${NEXT_TAG}"
          '''
        }
      }

      post {
        always {
          sh 'rm -f gpgkey.asc'
          publishHTML target: [
            allowMissing: false,
            alwaysLinkToLastBuild: false,
            keepAll: true,
            reportDir: 'coverage/',
            reportFiles: '*.html',
            reportName: 'Code Coverage'
          ]
          archiveArtifacts 'dist/*'
        }
      }
    }

    stage('Release') {
      when {
        beforeAgent true
        branch MAIN_BRANCH
      }

      agent {
        node {
          label 'ec2-fleet'
          customWorkspace "${PROJECT_NAME}-${BUILD_NUMBER}"
        }
      }

      environment {
        GIT_BRANCH = "${CURRENT_BRANCH}"
        GIT_AUTHOR = "${GIT_AUTHOR}"
        GIT_REPO = "${GIT_REPO}"
      }

      steps {
        script {
          withCredentials([
            string(credentialsId: 'logdna-gpg-key', variable: 'GPG_KEY'),
            string(credentialsId: 'github-api-token', variable: 'GITHUB_TOKEN')
          ]) {
            sh '''
              set +x
              git checkout -b ${GIT_BRANCH} origin/${GIT_BRANCH}
              git config user.name "LogDNA Bot"
              git config user.email "bot@logdna.com"
              git config remote.origin.url "https://${GIT_AUTHOR}:${GITHUB_TOKEN}@github.com/${GIT_REPO}"

              git fetch --tags
              export NEXT_TAG=$(make version-next)
              echo "Creating release for ${NEXT_TAG}"

              git tag ${NEXT_TAG}
              git push origin ${NEXT_TAG}
              echo "$GPG_KEY" > gpgkey.asc
              make release
            '''
          }
        }
      }

      post {
        always {
          sh 'rm -f gpgkey.asc'
        }
      }
    }
  }
}
