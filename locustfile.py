from locust import HttpUser, between, task


class WebsiteUser(HttpUser):
    wait_time = between(5, 15)

    @task
    def index1(self):
        self.client.post("/SetStr",{"Description":"张三录入礼品信息","Creator":"张三","AvailableTimes":"100","List":"士兵,4,炮车,8","ValidPeriod":"2021-09-02 00:00:00","CodeType":"3","CanGetUser":"212"})

    @task
    def index2(self):
        self.client.post("/GetStr",{"key":"3cb604a9"})

    @task
    def index3(self):
        self.client.post("/UpdateStr",{"key":"3cb604a9","username":"王武"})