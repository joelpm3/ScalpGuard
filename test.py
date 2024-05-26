import requests

# List of test user agent strings corresponding to various SEO bots
user_agents = [
    "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
    "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
    "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
    "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
    "rogerbot/1.0 (http://moz.com/help/guide/rogerbot-crawler)"
]

def test_user_agents(user_agents):
    for agent in user_agents:
        headers = {'User-Agent': agent}
        response = requests.get("https://127.0.0.1/", headers=headers, verify=False)
        print(response.text)
        print(f"Testing with User-Agent: {agent}")
        print(f"Response Status Code: {response.status_code}")
        print(f"Response Body: {response.text}\n")

test_user_agents(user_agents)