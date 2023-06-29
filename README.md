Cizzors

Cizzors
=======

Cizzors is a URL shortening service developed as the capstone project for the AltSchool Africa School of Back End Engineering. The project aims to disrupt the URL shortening industry by providing a simple tool to shorten URLs and customize them for branding purposes. Additionally, Cizzors offers QR code generation and basic analytics for tracking link performance.

Live Link
---------

The live version of Cizzors can be accessed at [https://cizzors.onrender.com](https://cizzors.onrender.com).

Running the App Locally using Docker
------------------------------------

To run the Cizzors app locally using Docker, follow these steps:

1.  Ensure that Docker is installed on your system. You can download and install Docker from the official website: [https://www.docker.com/get-started](https://www.docker.com/get-started).
2.  Clone the Cizzors repository from GitHub using the following command:

    git clone https://github.com/your-username/cizzors.git

3.  Navigate to the cloned repository directory:

    cd cizzors

4.  Build the Docker image using the provided Dockerfile:

    docker build -t cizzors-app .

5.  Once the image is built successfully, you can run the Docker container:

    docker run -p 8080:8080 cizzors-app

The Cizzors app should now be running locally on port 8080. You can access it in your browser by visiting [http://localhost:8080](http://localhost:8080).

Please note that the above steps assume that you have cloned the Cizzors repository and have Docker installed on your machine. Adjust the commands accordingly based on your environment and repository configuration.

This project was developed as the capstone project for the AltSchool Africa School of Back End Engineering.