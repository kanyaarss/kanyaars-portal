// Main JavaScript for Kanyaars Portal
// Enhanced with patterns from modern React applications

document.addEventListener('DOMContentLoaded', function() {
    console.log('Kanyaars Portal loaded');
    
    // Initialize all interactive features
    initLoader();
    initMobileMenu();
    initNavigationButtons();
    initHotspots();
    initRoomDimensions();
    initKeyboardNavigation();
    initMobileTitlePosition();
});

// ==================== LOADER/SPLASH SCREEN ====================
// Enhanced loader with smooth transitions (inspired by route-1.js)
function initLoader() {
    const loader = document.querySelector('.bg-white.fixed.top-0.left-0');
    const blurredContent = document.querySelector('[style*="filter:blur"]');
    
    if (loader) {
        // Use requestAnimationFrame for smoother animations
        setTimeout(() => {
            requestAnimationFrame(() => {
                loader.style.transition = 'opacity 1s ease-out, visibility 1s ease-out';
                loader.style.opacity = '0';
                loader.style.visibility = 'hidden';
                
                // Remove blur from content with smooth transition
                if (blurredContent) {
                    blurredContent.style.transition = 'filter 0.7s ease-out';
                    blurredContent.style.filter = 'blur(0px)';
                }
            });
        }, 1500);
    }
}

// ==================== MOBILE MENU ====================
function initMobileMenu() {
    // More specific selector for mobile menu button
    const mobileMenuLi = document.querySelector('li.md\\:hidden');
    const menuButton = mobileMenuLi ? mobileMenuLi.querySelector('button') : null;
    const menuDialog = document.querySelector('[role="dialog"][aria-modal="true"]');
    
    let currentMenuIndex = 0;
    const menuPanels = menuDialog ? menuDialog.querySelectorAll('.absolute.left-0.top-0') : [];
    
    // Toggle mobile menu
    if (menuButton && menuDialog) {
        menuButton.addEventListener('click', () => {
            const isOpen = !menuDialog.classList.contains('invisible');
            
            if (isOpen) {
                closeMenu();
            } else {
                openMenu(0);
            }
        });
    }
    
    function openMenu(panelIndex = 0) {
        if (!menuDialog) return;
        
        menuDialog.classList.remove('invisible', '-translate-x-full');
        menuDialog.classList.add('translate-x-0');
        
        // Apply background blur
        applyBackgroundBlur(true);
        
        // Show specific panel
        if (menuPanels[panelIndex]) {
            showPanel(panelIndex);
        }
    }
    
    function closeMenu() {
        if (!menuDialog) return;
        
        menuDialog.classList.add('invisible', '-translate-x-full');
        menuDialog.classList.remove('translate-x-0');
        
        // Remove background blur
        applyBackgroundBlur(false);
        
        // Hide all panels and reset items
        menuPanels.forEach(panel => {
            panel.classList.add('invisible', 'opacity-0');
            panel.classList.remove('visible', 'opacity-100');
            
            const items = panel.querySelectorAll('li');
            items.forEach(item => {
                item.classList.add('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                item.classList.remove('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
            });
        });
    }
    
    function showPanel(index) {
        // Hide all panels first and reset their items
        menuPanels.forEach((panel, i) => {
            if (i !== index) {
                panel.classList.add('invisible', 'opacity-0');
                panel.classList.remove('visible', 'opacity-100');
                
                // Reset all items in hidden panels
                const items = panel.querySelectorAll('li');
                items.forEach(item => {
                    item.classList.add('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                    item.classList.remove('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                });
            }
        });
        
        // Show selected panel
        const panel = menuPanels[index];
        if (panel) {
            panel.classList.remove('invisible', 'opacity-0');
            panel.classList.add('visible', 'opacity-100');
            
            // For panel 0 (mobile menu), only animate first UL, hide the mobile-explore-menu
            if (index === 0) {
                const firstUl = panel.querySelector('ul:not(.mobile-explore-menu)');
                const exploreUl = panel.querySelector('.mobile-explore-menu');
                
                // Hide explore menu items in panel 0
                if (exploreUl) {
                    const exploreItems = exploreUl.querySelectorAll('li');
                    exploreItems.forEach(item => {
                        item.classList.add('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                        item.classList.remove('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                    });
                }
                
                // Only animate first UL items
                if (firstUl) {
                    const items = firstUl.querySelectorAll('li');
                    items.forEach((item, i) => {
                        setTimeout(() => {
                            item.classList.remove('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                            item.classList.add('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                        }, i * 50);
                    });
                }
            } else {
                // For other panels, animate all items normally
                const items = panel.querySelectorAll('li');
                items.forEach((item, i) => {
                    setTimeout(() => {
                        item.classList.remove('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                        item.classList.add('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                    }, i * 50);
                });
            }
        }
    }
}

// ==================== NAVIGATION BUTTONS ====================
function initNavigationButtons() {
    // Use more specific selectors to target the correct buttons
    const navList = document.querySelector('.header-nav-list');
    const aboutButton = navList ? navList.querySelectorAll('button')[0] : null; // About Kay
    const projectsButton = navList ? navList.querySelectorAll('button')[1] : null; // Projects
    const exploreButton = navList ? navList.querySelectorAll('button')[2] : null; // Explore (fixed from index 3 to 2)
    const registeredButton = navList ? navList.querySelectorAll('button')[3] : null; // ® (fixed from index 4 to 3)
    
    const menuDialog = document.querySelector('[role="dialog"][aria-modal="true"]');
    const rightDialog = document.querySelectorAll('[role="dialog"][aria-modal="true"]')[1];
    
    // About Kay button - show first panel
    if (aboutButton) {
        aboutButton.addEventListener('click', () => {
            openLeftMenu(0);
        });
    }
    
    // Projects button - show second panel
    if (projectsButton) {
        projectsButton.addEventListener('click', () => {
            openLeftMenu(1);
        });
    }
    
    // Explore button - show fourth panel
    if (exploreButton) {
        exploreButton.addEventListener('click', () => {
            openLeftMenu(3);
        });
    }
    
    // Registered button - show right menu
    if (registeredButton) {
        registeredButton.addEventListener('click', () => {
            toggleRightMenu();
        });
    }
    
    function openLeftMenu(panelIndex) {
        if (!menuDialog) return;
        
        const menuPanels = menuDialog.querySelectorAll('.absolute.left-0.top-0');
        
        menuDialog.classList.remove('invisible', '-translate-x-full');
        menuDialog.classList.add('translate-x-0');
        
        // Apply background blur effect (from route-1.js)
        applyBackgroundBlur(true);
        
        // Hide all panels and reset their items
        menuPanels.forEach((panel, i) => {
            if (i !== panelIndex) {
                panel.classList.add('invisible', 'opacity-0');
                panel.classList.remove('visible', 'opacity-100');
                
                // Reset all items in hidden panels
                const items = panel.querySelectorAll('li');
                items.forEach(item => {
                    item.classList.add('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                    item.classList.remove('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                });
            }
        });
        
        // Show selected panel
        const panel = menuPanels[panelIndex];
        if (panel) {
            panel.classList.remove('invisible', 'opacity-0');
            panel.classList.add('visible', 'opacity-100');
            
            // Special handling for panel 0 (mobile menu) - hide explore items
            if (panelIndex === 0) {
                const firstUl = panel.querySelector('ul:not(.mobile-explore-menu)');
                const exploreUl = panel.querySelector('.mobile-explore-menu');
                
                // Hide explore menu items
                if (exploreUl) {
                    const exploreItems = exploreUl.querySelectorAll('li');
                    exploreItems.forEach(item => {
                        item.classList.add('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                        item.classList.remove('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                    });
                }
                
                // Only animate first UL items (About Kay, Projects, Tour)
                if (firstUl) {
                    const items = firstUl.querySelectorAll('li');
                    items.forEach((item, i) => {
                        setTimeout(() => {
                            item.classList.remove('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                            item.classList.add('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                        }, i * 50);
                    });
                }
            } else {
                // For other panels, animate all items normally
                const items = panel.querySelectorAll('li');
                items.forEach((item, i) => {
                    setTimeout(() => {
                        item.classList.remove('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                        item.classList.add('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                    }, i * 50);
                });
            }
        }
    }
    
    function toggleRightMenu() {
        if (!rightDialog) return;
        
        const isOpen = !rightDialog.classList.contains('invisible');
        
        if (isOpen) {
            rightDialog.classList.add('invisible', 'translate-x-full');
            rightDialog.classList.remove('translate-x-0');
            
            // Remove background blur
            applyBackgroundBlur(false);
            
            const panel = rightDialog.querySelector('.absolute.left-0.top-0');
            if (panel) {
                panel.classList.add('invisible', 'opacity-0');
                panel.classList.remove('visible', 'opacity-100');
                
                // Reset items
                const items = panel.querySelectorAll('li');
                items.forEach(item => {
                    item.classList.add('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                    item.classList.remove('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                });
            }
        } else {
            rightDialog.classList.remove('invisible', 'translate-x-full');
            rightDialog.classList.add('translate-x-0');
            
            // Apply background blur
            applyBackgroundBlur(true);
            
            const panel = rightDialog.querySelector('.absolute.left-0.top-0');
            if (panel) {
                panel.classList.remove('invisible', 'opacity-0');
                panel.classList.add('visible', 'opacity-100');
                
                // Animate items
                const items = panel.querySelectorAll('li');
                items.forEach((item, i) => {
                    setTimeout(() => {
                        item.classList.remove('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
                        item.classList.add('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
                    }, i * 50);
                });
            }
        }
    }
}

// ==================== INTERACTIVE HOTSPOTS ====================
// Enhanced hotspot interactions with touch support (from route-2.js beacon components)
function initHotspots() {
    const hotspots = document.querySelectorAll('#room a');
    
    hotspots.forEach(hotspot => {
        const button = hotspot.querySelector('span');
        const text = button ? button.querySelector('span:last-child') : null;
        
        if (button && text) {
            let isExpanded = false;
            
            // Calculate text width once for better performance
            const textWidth = text.scrollWidth;
            
            // Mouse hover effect - expand button and show text
            hotspot.addEventListener('mouseenter', () => {
                if (!isTouchDevice()) {
                    expandHotspot(button, text, textWidth);
                }
            });
            
            hotspot.addEventListener('mouseleave', () => {
                if (!isTouchDevice()) {
                    collapseHotspot(button, text);
                    isExpanded = false;
                }
            });
            
            // Touch/click support for mobile devices
            hotspot.addEventListener('click', (e) => {
                if (isTouchDevice() && !isExpanded) {
                    e.preventDefault();
                    expandHotspot(button, text, textWidth);
                    isExpanded = true;
                    
                    // Auto-collapse after 3 seconds
                    setTimeout(() => {
                        collapseHotspot(button, text);
                        isExpanded = false;
                    }, 3000);
                }
            });
        }
    });
}

// Helper function to expand hotspot
function expandHotspot(button, text, textWidth) {
    button.style.transition = 'width 0.5s ease-in-out';
    button.style.width = `${textWidth}px`;
    text.style.transition = 'opacity 0.5s ease-in-out';
    text.style.opacity = '1';
}

// Helper function to collapse hotspot
function collapseHotspot(button, text) {
    button.style.transition = 'width 0.5s ease-in-out';
    button.style.width = '24px';
    text.style.transition = 'opacity 0.5s ease-in-out';
    text.style.opacity = '0';
}

// Detect if device supports touch
function isTouchDevice() {
    return ('ontouchstart' in window) || 
           (navigator.maxTouchPoints > 0) || 
           (navigator.msMaxTouchPoints > 0);
}

// ==================== ROOM DIMENSIONS ====================
// Enhanced with aspect ratio handling and debounced resize (from route-2.js)
function initRoomDimensions() {
    const room = document.getElementById('room');
    const roomWrapper = document.querySelector('.room-wrapper');
    const DESKTOP_ASPECT_RATIO = 1.6; // 16:10 aspect ratio for desktop
    const MOBILE_ASPECT_RATIO = 0.75; // 3:4 aspect ratio for mobile (portrait)
    
    if (room) {
        function updateRoomSize() {
            const windowWidth = window.innerWidth;
            const windowHeight = window.innerHeight;
            const isMobile = windowWidth < 768; // md breakpoint
            const windowAspectRatio = windowWidth / windowHeight;
            
            let roomWidth, roomHeight;
            
            if (isMobile) {
                // Mobile: Use different aspect ratio and sizing strategy
                // Make room fit viewport height and allow horizontal scroll
                roomHeight = windowHeight;
                roomWidth = roomHeight * 1.5; // Wider for horizontal scrolling
                
                // Ensure minimum width for proper display
                if (roomWidth < windowWidth * 1.2) {
                    roomWidth = windowWidth * 1.5;
                    roomHeight = roomWidth / 1.5;
                }
            } else {
                // Desktop: Original logic
                const ASPECT_RATIO = DESKTOP_ASPECT_RATIO;
                if (windowAspectRatio < ASPECT_RATIO) {
                    roomHeight = windowHeight;
                    roomWidth = windowHeight * ASPECT_RATIO;
                } else {
                    roomWidth = windowWidth;
                    roomHeight = windowWidth / ASPECT_RATIO;
                }
            }
            
            room.style.width = `${roomWidth}px`;
            room.style.height = `${roomHeight}px`;
            
            // Center the room horizontally if it's wider than viewport
            if (roomWrapper && roomWidth > windowWidth) {
                const scrollLeft = (roomWidth - windowWidth) / 2;
                roomWrapper.scrollTo({
                    left: scrollLeft,
                    behavior: 'smooth'
                });
            }
        }
        
        // Debounce resize handler for better performance
        let resizeTimeout;
        function debouncedResize() {
            clearTimeout(resizeTimeout);
            resizeTimeout = setTimeout(updateRoomSize, 100);
        }
        
        updateRoomSize();
        window.addEventListener('resize', debouncedResize);
        
        // Add horizontal scroll support with mouse wheel (from route-2.js)
        if (roomWrapper) {
            initHorizontalScroll(roomWrapper);
        }
    }
}

// Enable horizontal scrolling with mouse wheel
function initHorizontalScroll(element) {
    element.addEventListener('wheel', (e) => {
        if (e.deltaY !== 0) {
            // Check if we're at scroll boundaries
            const atLeftEdge = element.scrollLeft === 0 && e.deltaY < 0;
            const atRightEdge = 
                element.scrollWidth - element.clientWidth - Math.round(element.scrollLeft) === 0 && 
                e.deltaY > 0;
            
            if (!atLeftEdge && !atRightEdge) {
                e.preventDefault();
                element.scrollTo({
                    left: element.scrollLeft + e.deltaY * 1.2
                });
            }
        }
    }, { passive: false });
}

// ==================== MOBILE TITLE POSITION ====================
// Keep title centered on mobile when scrolling (from route-2.js)
function initMobileTitlePosition() {
    const title = document.querySelector('h1.left-1\\/2');
    const roomWrapper = document.querySelector('.room-wrapper');
    
    if (!title || !roomWrapper) return;
    
    // Only apply on mobile
    const isMobile = window.innerWidth < 768;
    
    if (isMobile) {
        // Update title position on scroll to keep it centered
        roomWrapper.addEventListener('scroll', () => {
            const scrollLeft = roomWrapper.scrollLeft;
            const viewportWidth = window.innerWidth;
            const centerOffset = scrollLeft + (viewportWidth / 2);
            
            // Keep title at center of viewport, not center of room
            title.style.left = `${centerOffset}px`;
            title.style.transform = 'translateX(-50%)';
        });
        
        // Initial position
        const initialScroll = roomWrapper.scrollLeft;
        const viewportWidth = window.innerWidth;
        const centerOffset = initialScroll + (viewportWidth / 2);
        title.style.left = `${centerOffset}px`;
        title.style.transform = 'translateX(-50%)';
    }
    
    // Re-initialize on resize
    window.addEventListener('resize', () => {
        const newIsMobile = window.innerWidth < 768;
        if (newIsMobile) {
            const scrollLeft = roomWrapper.scrollLeft;
            const viewportWidth = window.innerWidth;
            const centerOffset = scrollLeft + (viewportWidth / 2);
            title.style.left = `${centerOffset}px`;
            title.style.transform = 'translateX(-50%)';
        } else {
            // Reset to original CSS on desktop
            title.style.left = '';
            title.style.transform = '';
        }
    });
}

// ==================== UTILITY FUNCTIONS ====================

// Utility function to make API calls
async function apiCall(endpoint, method = 'GET', data = null) {
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    };

    if (data) {
        options.body = JSON.stringify(data);
    }

    try {
        const response = await fetch(`/api/v1${endpoint}`, options);
        return await response.json();
    } catch (error) {
        console.error('API call failed:', error);
        throw error;
    }
}

// Utility function to show notifications
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    document.body.appendChild(notification);

    setTimeout(() => {
        notification.remove();
    }, 3000);
}

// ==================== BACKGROUND BLUR ====================
// Apply/remove background blur effect when menus are opened (from route-1.js)
function applyBackgroundBlur(shouldBlur) {
    const mainContent = document.getElementById('mainContent');
    const roomSection = document.querySelector('section');
    
    if (shouldBlur) {
        // Apply blur effect with smooth transition
        if (mainContent) {
            mainContent.style.transition = 'filter 0.5s ease-in-out';
            mainContent.style.filter = 'blur(8px)';
        }
        if (roomSection) {
            roomSection.style.transition = 'filter 0.5s ease-in-out';
            roomSection.style.filter = 'blur(8px)';
        }
    } else {
        // Remove blur effect
        if (mainContent) {
            mainContent.style.transition = 'filter 0.5s ease-in-out';
            mainContent.style.filter = 'blur(0px)';
        }
        if (roomSection) {
            roomSection.style.transition = 'filter 0.5s ease-in-out';
            roomSection.style.filter = 'blur(0px)';
        }
    }
}

// ==================== KEYBOARD NAVIGATION ====================
// Enhanced keyboard support (inspired by moda.js)
function initKeyboardNavigation() {
    document.addEventListener('keydown', (e) => {
        // Close menus on Escape key
        if (e.key === 'Escape') {
            const menuDialog = document.querySelector('[role="dialog"][aria-modal="true"]');
            const rightDialog = document.querySelectorAll('[role="dialog"][aria-modal="true"]')[1];
            
            if (menuDialog && !menuDialog.classList.contains('invisible')) {
                closeLeftMenu();
            }
            
            if (rightDialog && !rightDialog.classList.contains('invisible')) {
                closeRightMenu();
            }
        }
    });
}

// ==================== CLICK AWAY DETECTION ====================
// Enhanced click-away with better event delegation (from useclickaway.js)
document.addEventListener('click', (e) => {
    const menuDialog = document.querySelector('[role="dialog"][aria-modal="true"]');
    const rightDialog = document.querySelectorAll('[role="dialog"][aria-modal="true"]')[1];
    const navButtons = document.querySelectorAll('.header-nav-list button');
    const mobileMenuButton = document.querySelector('.md\\:hidden button');
    
    // Check if click is on navigation buttons or mobile menu button
    let isNavButton = false;
    navButtons.forEach(btn => {
        if (btn.contains(e.target)) {
            isNavButton = true;
        }
    });
    
    if (mobileMenuButton && mobileMenuButton.contains(e.target)) {
        isNavButton = true;
    }
    
    // Close menus if clicking outside
    if (!isNavButton) {
        if (menuDialog && !menuDialog.contains(e.target) && !menuDialog.classList.contains('invisible')) {
            closeLeftMenu();
        }
        
        if (rightDialog && !rightDialog.contains(e.target) && !rightDialog.classList.contains('invisible')) {
            closeRightMenu();
        }
    }
});

// Helper function to close left menu with animation
function closeLeftMenu() {
    const menuDialog = document.querySelector('[role="dialog"][aria-modal="true"]');
    if (!menuDialog) return;
    
    const menuPanels = menuDialog.querySelectorAll('.absolute.left-0.top-0');
    
    menuDialog.classList.add('invisible', '-translate-x-full');
    menuDialog.classList.remove('translate-x-0');
    
    // Remove background blur
    applyBackgroundBlur(false);
    
    // Hide all panels and reset items
    menuPanels.forEach(panel => {
        panel.classList.add('invisible', 'opacity-0');
        panel.classList.remove('visible', 'opacity-100');
        
        const items = panel.querySelectorAll('li');
        items.forEach(item => {
            item.classList.add('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
            item.classList.remove('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
        });
    });
}

// Helper function to close right menu with animation
function closeRightMenu() {
    const rightDialog = document.querySelectorAll('[role="dialog"][aria-modal="true"]')[1];
    if (!rightDialog) return;
    
    rightDialog.classList.add('invisible', 'translate-x-full');
    rightDialog.classList.remove('translate-x-0');
    
    // Remove background blur
    applyBackgroundBlur(false);
    
    const panel = rightDialog.querySelector('.absolute.left-0.top-0');
    if (panel) {
        panel.classList.add('invisible', 'opacity-0');
        panel.classList.remove('visible', 'opacity-100');
        
        const items = panel.querySelectorAll('li');
        items.forEach(item => {
            item.classList.add('opacity-0', 'invisible', '-translate-x-[10px]', 'pointer-events-none');
            item.classList.remove('opacity-100', 'visible', 'translate-x-0', 'pointer-events-auto');
        });
    }
}
