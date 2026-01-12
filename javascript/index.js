require('dotenv').config();

const apiToken = process.env.TOKEN;

async function getAllTasks(teamId, token) {
    let allTasks = [];
    let page = 0;
    let hasMore = true;
    const BATCH = 10; 
    const WAIT = 1000; 

    console.log(`starting fetch for Team ${teamId}...`);
    const startTime = Date.now();

    while (hasMore) {
        const pageBatch = Array.from({ length: BATCH }, (_, i) => page + i);
        
        try {
            const results = await Promise.all(
                pageBatch.map(p => fetchPage(teamId, token, p))
            );

            for (const tasks of results) {
                if (tasks && tasks.length > 0) {
                    allTasks.push(...tasks);
                } else {
                    hasMore = false; 
                }
            }

            page += BATCH;
            
            if (hasMore) await new Promise(res => setTimeout(res, WAIT));

        } catch (error) {
            console.error("batch failed:", error.message);
            hasMore = false;
        }
    }

    const performance = calculateMetrics(allTasks.length, page, startTime);
    
    console.log(`fetched ${performance.tasks} tasks in ${performance.duration} seconds. RMP: ${performance.RPM}`)
    
    return allTasks;
}

function calculateMetrics(totalTasks, totalPages, startTime) {
    const durationMs = Date.now() - startTime;
    const durationSec = durationMs / 1000;
    const durationMin = durationMs / 60000;

    return {
        "tasks": totalTasks,
        "duration": durationSec.toFixed(2),
        "RPM": (totalPages / durationMin).toFixed(2),
    };
}

async function fetchPage(teamId, token, page) {
    const response = await fetch(
        `https://api.clickup.com/api/v2/team/${teamId}/task?page=${page}&include_closed=true&subtasks=true`,
        {
            headers: { 
                'Authorization': `Bearer ${token}`, 
                'Content-Type': 'application/json' 
            }
        }
    );
    
    if (response.status === 429) throw new Error("429: rate limit");
    if (!response.ok) throw new Error(`${response.status}: ${response.statusText}`);
    
    const data = await response.json();
    return data.tasks || [];
}

getAllTasks("36226098", apiToken);