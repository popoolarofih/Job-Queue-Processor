import React, { useState, useEffect, useCallback } from 'react';
import './App.css';

function App() {
  const [jobName, setJobName] = useState('');
  const [submitMessage, setSubmitMessage] = useState('');
  const [queuesData, setQueuesData] = useState([]);
  const [isLoadingStatus, setIsLoadingStatus] = useState(false);
  const [errorStatus, setErrorStatus] = useState('');

  const fetchJobStatus = useCallback(async () => {
    setIsLoadingStatus(true);
    setErrorStatus('');
    try {
      const response = await fetch('/api/jobs/status');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setQueuesData(data || []); // Ensure data is an array
    } catch (error) {
      console.error("Failed to fetch job statuses:", error);
      setErrorStatus(`Failed to load job statuses: ${error.message}`);
      setQueuesData([]); // Clear data on error
    } finally {
      setIsLoadingStatus(false);
    }
  }, []);

  useEffect(() => {
    fetchJobStatus(); // Initial fetch
    const intervalId = setInterval(fetchJobStatus, 5000); // Poll every 5 seconds
    return () => clearInterval(intervalId); // Cleanup on unmount
  }, [fetchJobStatus]);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setSubmitMessage('');

    if (!jobName.trim()) {
      setSubmitMessage('Please enter a name for the job.');
      return;
    }

    try {
      const response = await fetch('/api/jobs', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          type: 'greeting:sayhello',
          payload: { name: jobName },
        }),
      });

      const data = await response.json();

      if (response.ok) {
        setSubmitMessage(`Job submitted successfully! ID: ${data.job_id}, Type: ${data.type}, Queue: ${data.queue}`);
        setJobName(''); // Clear input
        fetchJobStatus(); // Refresh status immediately after submission
      } else {
        setSubmitMessage(`Error submitting job: ${data.error || 'Unknown error'}`);
      }
    } catch (error) {
      setSubmitMessage(`Network error: ${error.message}`);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Job Queue Dashboard</h1>
      </header>
      <main>
        <section className="job-submission-form">
          <h2>Submit New Greeting Job</h2>
          <form onSubmit={handleSubmit}>
            <div>
              <label htmlFor="jobName">Name:</label>
              <input
                type="text"
                id="jobName"
                value={jobName}
                onChange={(e) => setJobName(e.target.value)}
                placeholder="Enter a name for the greeting"
              />
            </div>
            <button type="submit">Submit Job</button>
          </form>
          {submitMessage && <p className={`message ${submitMessage.startsWith('Error') || submitMessage.startsWith('Network error') ? 'error' : 'success'}`}>{submitMessage}</p>}
        </section>

        <section className="job-status-display">
          <h2>Current Job Statuses</h2>
          <button onClick={fetchJobStatus} disabled={isLoadingStatus}>
            {isLoadingStatus ? 'Refreshing...' : 'Refresh Status'}
          </button>
          {errorStatus && <p className="message error">{errorStatus}</p>}

          {queuesData.length === 0 && !isLoadingStatus && !errorStatus && <p>No queue data available or all queues are empty.</p>}

          {queuesData.map((queue) => (
            <div key={queue.name} className="queue-info">
              <h3>Queue: {queue.name}</h3>
              <p><strong>Size:</strong> {queue.size} | <strong>Active:</strong> {queue.active} | <strong>Pending:</strong> {queue.pending} | <strong>Completed:</strong> {queue.completed} | <strong>Retry:</strong> {queue.retry} | <strong>Archived:</strong> {queue.archived}</p>
              <p><strong>Processed:</strong> {queue.processed} | <strong>Failed:</strong> {queue.failed}</p>
              <h4>Tasks (Recent 5 per category):</h4>
              {queue.tasks && queue.tasks.length > 0 ? (
                <ul className="job-list">
                  {queue.tasks.map((task) => (
                    <li key={task.id} className="job-item">
                      <p><strong>ID:</strong> {task.id}</p>
                      <p><strong>Type:</strong> {task.type}</p>
                      <p><strong>Payload:</strong> {task.payload}</p>
                      <p><strong>State:</strong> <span className={`task-state ${task.state}`}>{task.state}</span></p>
                      {task.last_err && <p><strong>Error:</strong> {task.last_err}</p>}
                    </li>
                  ))}
                </ul>
              ) : (
                <p>No tasks found in the listed categories for this queue.</p>
              )}
            </div>
          ))}
        </section>
      </main>
    </div>
  );
}

export default App;
