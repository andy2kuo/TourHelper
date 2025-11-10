-- Create tours table
CREATE TABLE IF NOT EXISTS tours (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(255) NOT NULL,
    country VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    duration INTEGER NOT NULL,
    season VARCHAR(50),
    budget VARCHAR(20),
    image_url TEXT,
    rating DECIMAL(3, 2) DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common queries
CREATE INDEX idx_tours_category ON tours(category);
CREATE INDEX idx_tours_country ON tours(country);
CREATE INDEX idx_tours_budget ON tours(budget);
CREATE INDEX idx_tours_season ON tours(season);
CREATE INDEX idx_tours_rating ON tours(rating);

-- Insert sample data
INSERT INTO tours (name, description, location, country, category, duration, season, budget, image_url, rating) VALUES
('Tokyo Cherry Blossom Tour', 'Experience the beautiful cherry blossoms in Tokyo with guided tours to famous spots', 'Tokyo', 'Japan', 'cultural', 5, 'spring', 'high', 'https://example.com/tokyo.jpg', 4.8),
('Bali Beach Paradise', 'Relax on pristine beaches and explore traditional Balinese culture', 'Bali', 'Indonesia', 'beach', 7, 'summer', 'medium', 'https://example.com/bali.jpg', 4.7),
('Swiss Alps Adventure', 'Hiking and skiing in the breathtaking Swiss Alps', 'Interlaken', 'Switzerland', 'mountain', 6, 'winter', 'high', 'https://example.com/swiss.jpg', 4.9),
('Bangkok Street Food Tour', 'Discover the amazing street food culture of Bangkok', 'Bangkok', 'Thailand', 'cultural', 3, 'all', 'low', 'https://example.com/bangkok.jpg', 4.6),
('Great Barrier Reef Diving', 'Explore the world''s largest coral reef system', 'Queensland', 'Australia', 'beach', 4, 'summer', 'high', 'https://example.com/reef.jpg', 4.9),
('Paris City Break', 'Visit iconic landmarks and enjoy French cuisine', 'Paris', 'France', 'city', 4, 'spring', 'medium', 'https://example.com/paris.jpg', 4.7);
