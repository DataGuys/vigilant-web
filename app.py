# app.py
from flask import Flask, request, jsonify, render_template_string
import subprocess, sqlite3, json, os

app = Flask(__name__)
DATABASE = 'scan_results.db'

def init_db():
    conn = sqlite3.connect(DATABASE)
    c = conn.cursor()
    c.execute('''
        CREATE TABLE IF NOT EXISTS scans (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            term TEXT,
            result TEXT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    ''')
    conn.commit()
    conn.close()

init_db()

@app.route('/')
def index():
    html = '''
    <h1>Vigilant Onion Darkweb Scanner</h1>
    <form action="/scan" method="post">
        <label>Enter search term (e.g., your company name):</label><br>
        <input type="text" name="term" required>
        <input type="submit" value="Scan">
    </form>
    <br>
    <a href="/results">View Scan Results</a>
    '''
    return render_template_string(html)

@app.route('/scan', methods=['POST'])
def scan():
    term = request.form.get('term')
    if not term:
        return "Search term missing", 400
    try:
        # Call observer.py with the --find flag
        cmd = ["python3", "observer.py", "--config", "config/config.yml", "--find", term]
        output = subprocess.check_output(cmd, stderr=subprocess.STDOUT, timeout=300)
        result = output.decode('utf-8')
    except subprocess.CalledProcessError as e:
        result = f"Error: {e.output.decode('utf-8')}"
    except Exception as ex:
        result = f"Exception: {str(ex)}"
    
    # Store result in SQLite database
    conn = sqlite3.connect(DATABASE)
    c = conn.cursor()
    c.execute("INSERT INTO scans (term, result) VALUES (?, ?)", (term, result))
    conn.commit()
    conn.close()
    return jsonify({"term": term, "result": result})

@app.route('/results')
def results():
    conn = sqlite3.connect(DATABASE)
    c = conn.cursor()
    c.execute("SELECT id, term, result, timestamp FROM scans ORDER BY timestamp DESC")
    rows = c.fetchall()
    conn.close()
    # Return JSON array of scan results
    results_list = [{"id": row[0], "term": row[1], "result": row[2], "timestamp": row[3]} for row in rows]
    return jsonify(results_list)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
