// Application state
const state = {
    token: null,
    username: null,
    platform: null,
    ws: null,
    map: null,
    markers: [],
    googleMapsApiKey: ''
};

// Initialize the app
document.addEventListener('DOMContentLoaded', () => {
    // Check if already logged in
    const token = localStorage.getItem('token');
    if (token) {
        state.token = token;
        state.username = localStorage.getItem('username');
        state.platform = localStorage.getItem('platform');
        showMainPage();
    }

    setupEventListeners();
});

function setupEventListeners() {
    // Login form
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }

    // Logout button
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', handleLogout);
    }

    // Search button
    const searchBtn = document.getElementById('search-btn');
    if (searchBtn) {
        searchBtn.addEventListener('click', handleSearch);
    }

    // Location search enter key
    const locationSearch = document.getElementById('location-search');
    if (locationSearch) {
        locationSearch.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                handleSearch();
            }
        });
    }

    // Get suggestions button
    const getSuggestionsBtn = document.getElementById('get-suggestions-btn');
    if (getSuggestionsBtn) {
        getSuggestionsBtn.addEventListener('click', handleGetSuggestions);
    }

    // Send message button
    const sendBtn = document.getElementById('send-btn');
    if (sendBtn) {
        sendBtn.addEventListener('click', handleSendMessage);
    }

    // Message input enter key
    const messageInput = document.getElementById('message-input');
    if (messageInput) {
        messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                handleSendMessage();
            }
        });
    }
}

async function handleLogin(e) {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const platform = document.getElementById('platform').value;
    const errorDiv = document.getElementById('login-error');

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password, platform })
        });

        if (!response.ok) {
            throw new Error('Login failed');
        }

        const data = await response.json();
        
        // Save to state and localStorage
        state.token = data.token;
        state.username = data.username;
        state.platform = data.platform;
        
        localStorage.setItem('token', data.token);
        localStorage.setItem('username', data.username);
        localStorage.setItem('platform', data.platform);

        showMainPage();
    } catch (error) {
        errorDiv.textContent = 'Invalid username or password';
    }
}

function handleLogout() {
    // Clear state
    state.token = null;
    state.username = null;
    state.platform = null;
    
    // Clear localStorage
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    localStorage.removeItem('platform');

    // Close WebSocket
    if (state.ws) {
        state.ws.close();
        state.ws = null;
    }

    // Show login page
    showLoginPage();
}

function showLoginPage() {
    document.getElementById('login-page').classList.add('active');
    document.getElementById('main-page').classList.remove('active');
}

async function showMainPage() {
    document.getElementById('login-page').classList.remove('active');
    document.getElementById('main-page').classList.add('active');

    // Update user info
    document.getElementById('user-name').textContent = state.username;
    document.getElementById('user-platform').textContent = state.platform;

    // Get config
    await loadConfig();

    // Initialize map
    initializeMap();

    // Connect WebSocket
    connectWebSocket();
}

async function loadConfig() {
    try {
        const response = await fetch('/api/config', {
            headers: {
                'Authorization': `Bearer ${state.token}`
            }
        });

        if (response.ok) {
            const config = await response.json();
            state.googleMapsApiKey = config.googleMapsApiKey;
        }
    } catch (error) {
        console.error('Failed to load config:', error);
    }
}

function initializeMap() {
    // Load Google Maps API dynamically
    if (!window.google) {
        const script = document.createElement('script');
        script.src = `https://maps.googleapis.com/maps/api/js?key=${state.googleMapsApiKey}&libraries=places&callback=initMap`;
        script.async = true;
        script.defer = true;
        document.head.appendChild(script);
    } else {
        initMap();
    }
}

window.initMap = function() {
    // Initialize map centered on Taiwan
    state.map = new google.maps.Map(document.getElementById('map'), {
        center: { lat: 23.6978, lng: 120.9605 },
        zoom: 8
    });

    // Add a marker
    addMarker({ lat: 23.6978, lng: 120.9605 }, 'Taiwan');
};

function addMarker(position, title) {
    const marker = new google.maps.Marker({
        position: position,
        map: state.map,
        title: title
    });

    state.markers.push(marker);
    return marker;
}

function clearMarkers() {
    state.markers.forEach(marker => marker.setMap(null));
    state.markers = [];
}

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?token=${state.token}`;

    state.ws = new WebSocket(wsUrl);

    state.ws.onopen = () => {
        console.log('WebSocket connected');
        updateWSStatus(true);
    };

    state.ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        handleWSMessage(message);
    };

    state.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        updateWSStatus(false);
    };

    state.ws.onclose = () => {
        console.log('WebSocket disconnected');
        updateWSStatus(false);
        
        // Reconnect after 3 seconds
        setTimeout(() => {
            if (state.token) {
                connectWebSocket();
            }
        }, 3000);
    };
}

function updateWSStatus(connected) {
    const indicator = document.getElementById('ws-indicator');
    const statusText = document.getElementById('ws-status-text');

    if (connected) {
        indicator.classList.add('connected');
        statusText.textContent = 'Connected';
    } else {
        indicator.classList.remove('connected');
        statusText.textContent = 'Disconnected';
    }
}

function handleWSMessage(message) {
    console.log('WebSocket message:', message);

    if (message.type === 'chat') {
        addChatMessage(message.payload);
    } else if (message.type === 'suggestion') {
        displaySuggestion(message.payload);
    }
}

function handleSearch() {
    const query = document.getElementById('location-search').value;
    if (!query) return;

    if (!window.google) {
        alert('Google Maps is loading, please wait...');
        return;
    }

    const service = new google.maps.places.PlacesService(state.map);
    const request = {
        query: query,
        fields: ['name', 'geometry']
    };

    service.findPlaceFromQuery(request, (results, status) => {
        if (status === google.maps.places.PlacesServiceStatus.OK && results) {
            clearMarkers();
            
            results.forEach(place => {
                if (place.geometry && place.geometry.location) {
                    addMarker(place.geometry.location, place.name);
                    state.map.setCenter(place.geometry.location);
                    state.map.setZoom(15);
                }
            });
        } else {
            alert('Location not found');
        }
    });
}

function handleGetSuggestions() {
    const suggestions = [
        {
            id: '1',
            name: 'Taipei 101',
            description: 'Iconic skyscraper and observation tower',
            location: { lat: 25.0340, lng: 121.5645 }
        },
        {
            id: '2',
            name: 'Sun Moon Lake',
            description: 'Beautiful alpine lake in central Taiwan',
            location: { lat: 23.8561, lng: 120.9142 }
        },
        {
            id: '3',
            name: 'Taroko Gorge',
            description: 'Spectacular marble canyon',
            location: { lat: 24.1939, lng: 121.4909 }
        }
    ];

    displaySuggestions(suggestions);

    // Send via WebSocket
    if (state.ws && state.ws.readyState === WebSocket.OPEN) {
        state.ws.send(JSON.stringify({
            type: 'suggestion',
            payload: { suggestions }
        }));
    }
}

function displaySuggestions(suggestions) {
    const listDiv = document.getElementById('suggestions-list');
    listDiv.innerHTML = '';

    suggestions.forEach(suggestion => {
        const item = document.createElement('div');
        item.className = 'suggestion-item';
        item.innerHTML = `
            <h4>${suggestion.name}</h4>
            <p>${suggestion.description}</p>
        `;
        item.onclick = () => {
            if (window.google && suggestion.location) {
                clearMarkers();
                addMarker(suggestion.location, suggestion.name);
                state.map.setCenter(suggestion.location);
                state.map.setZoom(12);
            }
        };
        listDiv.appendChild(item);
    });
}

function displaySuggestion(payload) {
    if (payload.suggestions) {
        displaySuggestions(payload.suggestions);
    }
}

function handleSendMessage() {
    const input = document.getElementById('message-input');
    const message = input.value.trim();

    if (!message) return;

    if (state.ws && state.ws.readyState === WebSocket.OPEN) {
        state.ws.send(JSON.stringify({
            type: 'chat',
            payload: {
                sender: state.username,
                message: message,
                platform: state.platform
            }
        }));

        input.value = '';
    } else {
        alert('WebSocket not connected');
    }
}

function addChatMessage(payload) {
    const messagesDiv = document.getElementById('messages');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message';

    const time = new Date().toLocaleTimeString();
    
    messageDiv.innerHTML = `
        <span class="sender">${payload.sender || 'Unknown'}:</span>
        <span class="content">${payload.message || ''}</span>
        <span class="time">${time}</span>
    `;

    messagesDiv.appendChild(messageDiv);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}
