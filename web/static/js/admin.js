// Admin JavaScript for Kanyaars Portal

document.addEventListener('DOMContentLoaded', function() {
    // Check if user is logged in
    const token = localStorage.getItem('token');
    if (!token && !window.location.pathname.includes('/admin/login')) {
        window.location.href = '/admin/login';
    }

    // Load dashboard data if on dashboard
    if (window.location.pathname === '/admin/') {
        loadDashboardData();
    }

    // Setup login form
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
});

// Handle admin login
async function handleLogin(e) {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const errorMessage = document.getElementById('error-message');

    try {
        const response = await fetch('/api/v1/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });

        const data = await response.json();

        if (data.success) {
            localStorage.setItem('token', data.data.token);
            localStorage.setItem('user', JSON.stringify(data.data.user));
            window.location.href = '/admin/';
        } else {
            errorMessage.textContent = data.error || 'Login failed';
            errorMessage.style.display = 'block';
        }
    } catch (error) {
        errorMessage.textContent = 'An error occurred. Please try again.';
        errorMessage.style.display = 'block';
        console.error('Login error:', error);
    }
}

// Load dashboard data
async function loadDashboardData() {
    try {
        const response = await fetch('/api/v1/projects', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });

        const data = await response.json();

        if (data.success && data.data && data.data.data) {
            const projects = data.data.data;
            
            // Update stats
            document.getElementById('total-projects').textContent = projects.length;
            const activeCount = projects.filter(p => p.status === 'active').length;
            const inactiveCount = projects.filter(p => p.status === 'inactive').length;
            
            document.getElementById('active-projects').textContent = activeCount;
            document.getElementById('inactive-projects').textContent = inactiveCount;

            // Update projects table
            const tableBody = document.getElementById('projects-table');
            tableBody.innerHTML = '';

            projects.forEach(project => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${project.name}</td>
                    <td>${project.slug}</td>
                    <td><span class="badge">${project.status}</span></td>
                    <td>
                        <a href="/admin/projects/${project.id}" class="btn btn-secondary">Edit</a>
                        <button onclick="deleteProject(${project.id})" class="btn btn-danger">Delete</button>
                    </td>
                `;
                tableBody.appendChild(row);
            });
        }
    } catch (error) {
        console.error('Failed to load dashboard data:', error);
    }
}

// Delete project
async function deleteProject(id) {
    if (!confirm('Are you sure you want to delete this project?')) {
        return;
    }

    try {
        const response = await fetch(`/api/v1/admin/projects/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });

        const data = await response.json();

        if (data.success) {
            loadDashboardData();
            showNotification('Project deleted successfully', 'success');
        } else {
            showNotification('Failed to delete project', 'error');
        }
    } catch (error) {
        console.error('Delete error:', error);
        showNotification('An error occurred', 'error');
    }
}

// Logout
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/admin/login';
}

// Show notification
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 1rem;
        background-color: ${type === 'success' ? '#2ecc71' : '#e74c3c'};
        color: white;
        border-radius: 4px;
        z-index: 1000;
    `;
    document.body.appendChild(notification);

    setTimeout(() => {
        notification.remove();
    }, 3000);
}
