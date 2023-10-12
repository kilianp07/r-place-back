# R-Place - Real-Time Pixel Board Application

R-Place is a real-time pixel board application that allows users to place colored pixels on a large shared canvas.

## Prerequisites

Before getting started, make sure you have the following installed on your system:

- [Docker](https://www.docker.com/get-started) (to run the application in a container)
- [Docker Compose](https://docs.docker.com/compose/install/) (for managing multiple containers)
- A modern web browser for viewing the user interface

## Server Setup

1. Clone this GitHub repository to your local machine.

```bash
git clone https://github.com/your-username/r-place
```

2. Navigate to the project directory.

```bash
cd r-place
```

3. Use Docker Compose to build and run the application and MySQL server containers.

```bash
docker-compose up
```

This will create and start the necessary containers to run the application, including the database.

4. Once the containers are ready, you can access the application by opening a web browser and going to [http://localhost:80](http://localhost:80).

## Using the Application

1. You will see an empty pixel board. You can click on any area of the board to place a colored pixel. 

2. Color updates are broadcast in real-time. You will see colored pixels placed by other users appearing instantly.

