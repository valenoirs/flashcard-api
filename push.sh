docker rmi gcr.io/valenoirs-production/flashcard-api:0.1.0
docker build -t gcr.io/valenoirs-production/flashcard-api:0.1.0 .
docker push gcr.io/valenoirs-production/flashcard-api:0.1.0
