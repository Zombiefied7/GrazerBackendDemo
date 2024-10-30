# README

[Project Design Document](./DESIGN.md)

## Setup Instructions

1. **Set Up Local Docker Registry**  
   Start a local Docker registry to store and access images.

   ```bash
   docker run -d -p 5000:5000 --name registry registry:2
   ```

2. **Build and Push Docker Images**

   ```bash
   cd grazer-user-service
   docker build -t localhost:5000/grazer-user-service:latest .
   docker push localhost:5000/grazer-user-service:latest

   cd ../grazer-preferences-service
   docker build -t localhost:5000/grazer-preferences-service:latest .
   docker push localhost:5000/grazer-preferences-service:latest

   cd ../grazer-liking-service
   docker build -t localhost:5000/grazer-liking-service:latest .
   docker push localhost:5000/grazer-liking-service:latest

   cd ../grazer-matching-service
   docker build -t localhost:5000/grazer-matching-service:latest .
   docker push localhost:5000/grazer-matching-service:latest
   ```

3. **Deploy Kubernetes Services and Pods**

   ```bash
   kubectl apply -f grazer-database-deployment.yaml
   kubectl apply -f grazer-user-service/grazer-user-service-deployment.yaml
   kubectl apply -f grazer-preferences-service/grazer-preferences-service-deployment.yaml
   kubectl apply -f grazer-liking-service/grazer-liking-service-deployment.yaml
   kubectl apply -f grazer-matching-service/grazer-matching-service-deployment.yaml
   kubectl apply -f rabbitmq-deployment.yaml
   ```

4. **Access RabbitMQ Management Interface**  
   Go to [http://localhost:31999/#/queues](http://localhost:31999/#/queues) and log in with:

   - **Username**: `guest`
   - **Password**: `guest`

5. **Verify Services are Running**
   ```bash
   kubectl get pods
   kubectl get services
   ```

## Important Notes

- **Service Startup Order**: The order in which services are started can affect functionality, especially with RabbitMQ, which must be fully initialized before dependent services (like the Matching and Liking services). This project is not optimized to handle startup dependencies automatically.
