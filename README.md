# <img src="./static/img.icons8.png" width="50" height="50"/> GoFit

#### Description
This is a personal project, it is a full fledged fitness tracker. There are 3 programs a user can select, fat loss, muscle gain and maintenance. Your calorie
deficit or surplus will be calculated based on the program you've selected and your weight, height, current weight and desired weight. It comes with a 
calorie calculator where you can enter the prompt in human readable text like "200 ml of whole milk" and it will output the nutrition of the input prompt.
You can even track your progress.
<br>
<br>

#### Motivation
Motivation for making this project was that I wanted an application to suggest the amount of calories something had and I wanted it to show me my deficit
if I am cutting or my surplus if I am bulking. I also was learning Go and HTMX and saw this as a great idea for a personal project.
<br>
<br>

## Interface
### Desktop

*Program selection*

<img src="./static/screenshots/homepage.png" /> 

*User logs*

<img src="./static/screenshots/user_logs.png" /> 

*User profile*

<img src="./static/screenshots/user_profile.png" /> 

*Daily input tracker*

<img src="./static/screenshots/daily_input_tracker.png" /> 

*Signup page*

<img src="./static/screenshots/signup.png" /> 

*Calorie calculator*

<img src="./static/screenshots/calorie_calc.png" /> 

### Mobile

<p>
    <img src="./static/screenshots/mobile_interface.png" width="49%" /> 
    <img src="./static/screenshots/mobile_homepage.png" width="49%"/>
</p>


## Run
#### Setting up environment
###### Installing Go
```bash
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gzmake build
export PATH=$PATH:/usr/local/go/bin
go version
```
###### Installing Docker && postgres
```bash
sudo apt update
sudo apt install curl
curl -fsSL https://get.docker.com/ | sh
docker pull postgres
docker run -d --name <name> -p 5432:5432 -e POSTGRES_PASSWORD=<password> postgres
docker exec -it <name> bash
psql -h localhost -U postgres
CREATE DATABASE fitness;
```
###### Create a .env file with these parameters
```env
DB_TYPE=postgres
DB_URL=postgres://postgres:<your-password>@172.17.0.2:5432/fitness?sslmode=disable
BASE_URL=http://localhost:5000
JWT_SECRET=<your-secret>
API_APP_KEY=<your-api-key>
API_ACCESS_POINT=<your-access-point>
API_APP_ID=<your-api-id>
```
###### For api use [edamam](https://developer.edamam.com/edamam-nutrition-api-demo)
#### Execution
```bash
docker start <name-of-postgres-continer>
make intall
make migrageUp
make queries
make build
make run
```

## Tech Stack
+ **Frontend**\
    ![Templ](https://img.shields.io/badge/Templ-%252300ADD8?style=for-the-badge&logo=go&logoColor=yellow&logoSize=auto&color=black)
    ![Htmx](https://img.shields.io/badge/htmx-36C?logo=htmx&logoColor=fff&style=for-the-badge&logoSize=auto)
    ![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?logo=javascript&logoColor=000&style=for-the-badge&logoSize=auto)
    ![TailwindCSS](https://img.shields.io/badge/tailwindcss-%2338B2AC.svg?style=for-the-badge&logo=tailwind-css&logoSize=auto&logoColor=white)
    ![DaisyUI](https://img.shields.io/badge/daisyui-5A0EF8?style=for-the-badge&logo=daisyui&logoSize=auto&logoColor=white)
+ **Backend**\
    ![Go](https://img.shields.io/badge/Golang-%252300ADD8?style=for-the-badge&logo=go&logoColor=white&logoSize=auto&color=blue)
+ **Database**\
    ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white&logoSize=auto)
+ **Libraries**\
    ![Chi](https://img.shields.io/badge/Chi-%252300ADD8?style=for-the-badge&logo=go&logoColor=green&logoSize=auto&color=white)
    ![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens&logoSize=auto)
+ **Tools**\
    ![Goose](https://img.shields.io/badge/Goose-%252300ADD8?style=for-the-badge&logo=go&logoColor=orange&logoSize=auto&color=white)
    ![Sqlc](https://img.shields.io/badge/Sqlc-%23316192.svg?style=for-the-badge&logo=postgresql&logoSize=auto&logoColor=yellow&color=black)
    ![Psql](https://img.shields.io/badge/psql-%23316192.svg?style=for-the-badge&logo=postgresql&logoSize=auto&logoColor=violet&color=black)
    ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoSize=auto&logoColor=white)
