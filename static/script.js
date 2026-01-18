const watchBtn = document.getElementById('watchBtn');
const playerContainer = document.getElementById('player');
let retryCount = 0;
const maxRetries = 3;

function loadRandomVideo() {
    retryCount = 0;
    watchBtn.disabled = true;
    watchBtn.textContent = 'Loading...';
    fetchAndDisplayVideo();
}

function fetchAndDisplayVideo() {
    fetch('/api/random')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            const videoURL = data.video_url;
            const title = data.title;
            const photographer = data.photographer;
            const duration = data.duration;
            
            playerContainer.innerHTML = `
                <div class="video-info">
                    <h3>${title}</h3>
                    <p>By <strong>${photographer}</strong> • ${duration}s</p>
                </div>
                <div class="video-wrapper">
                    <video controls autoplay>
                        <source src="${videoURL}" type="video/mp4">
                        Your browser does not support the video tag.
                    </video>
                </div>
            `;
            
            watchBtn.disabled = false;
            watchBtn.textContent = 'Watch Random Video';
        })
        .catch(error => {
            console.error('Error fetching random video:', error);
            retryCount++;
            
            if (retryCount < maxRetries) {
                console.log(`Error loading video. Retrying... (${retryCount}/${maxRetries})`);
                playerContainer.innerHTML = `<p>Loading... (Attempt ${retryCount + 1}/${maxRetries + 1})</p>`;
                setTimeout(fetchAndDisplayVideo, 2000);
            } else {
                playerContainer.innerHTML = `
                    <div style="padding: 20px; text-align: center;">
                        <p style="color: #e94560; font-weight: bold;">Could not load video. Please try again.</p>
                        <p style="color: #999; font-size: 0.9em;">Error: ${error.message}</p>
                    </div>
                `;
                watchBtn.disabled = false;
                watchBtn.textContent = 'Watch Random Video';
            }
        });
}

watchBtn.addEventListener('click', loadRandomVideo);