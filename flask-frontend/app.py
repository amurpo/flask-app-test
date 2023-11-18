from flask import Flask, render_template
import requests

app = Flask(__name__)

# Replace this URL with your Go GraphQL server endpoint
GRAPHQL_URL = "http://localhost:8000/graphql"

def fetch_images():
    query = '''
    {
        images {
            ID
            Link
        }
    }
    '''

    response = requests.post(GRAPHQL_URL, json={'query': query})
    data = response.json()
    return data['data']['images']

@app.route('/')
def index():
    images = fetch_images()

    return render_template('index.html', images=images)

if __name__ == '__main__':
    app.run(debug=True)

