#!/bin/bash

# Product Service Database Migration Script
# This script runs all migration files for the product service

set -e

# Configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-root}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-edot_product}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if MySQL is running
check_mysql() {
    print_info "Checking MySQL connection..."
    if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1;" > /dev/null 2>&1; then
        print_error "Cannot connect to MySQL. Please check your database configuration."
        exit 1
    fi
    print_success "MySQL connection established"
}

# Create database if it doesn't exist
create_database() {
    print_info "Creating database '$DB_NAME' if it doesn't exist..."
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;" 2>/dev/null
    print_success "Database '$DB_NAME' is ready"
}

# Create migration tracking table
create_migration_table() {
    print_info "Creating migration tracking table..."
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" << 'EOF'
CREATE TABLE IF NOT EXISTS migration_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_migration_filename (filename)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
EOF
    print_success "Migration tracking table created"
}

# Run migration file
run_migration() {
    local file="$1"
    local filename=$(basename "$file")
    
    # Check if migration already executed
    local count=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -se "SELECT COUNT(*) FROM migration_history WHERE filename = '$filename';" 2>/dev/null || echo "0")
    
    if [ "$count" -gt 0 ]; then
        print_warning "Migration '$filename' already executed, skipping..."
        return 0
    fi
    
    print_info "Executing migration: $filename"
    
    # Execute the migration file
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$file"; then
        # Record successful migration
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "INSERT INTO migration_history (filename) VALUES ('$filename');" 2>/dev/null
        print_success "Migration '$filename' completed successfully"
    else
        print_error "Migration '$filename' failed"
        exit 1
    fi
}

# Main execution
main() {
    print_info "Starting Product Service Database Migration"
    print_info "Database: $DB_NAME on $DB_HOST:$DB_PORT"
    
    # Get script directory
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    
    # Check MySQL connection
    check_mysql
    
    # Create database
    create_database
    
    # Create migration tracking table
    create_migration_table
    
    # Find and execute migration files in order
    print_info "Looking for migration files in: $SCRIPT_DIR"
    
    if [ ! -d "$SCRIPT_DIR" ]; then
        print_error "Migration directory not found: $SCRIPT_DIR"
        exit 1
    fi
    
    # Execute migrations in order
    for migration_file in "$SCRIPT_DIR"/*.sql; do
        if [ -f "$migration_file" ]; then
            run_migration "$migration_file"
        fi
    done
    
    print_success "All migrations completed successfully!"
    
    # Show summary
    print_info "Migration Summary:"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SELECT filename, executed_at FROM migration_history ORDER BY executed_at;" 2>/dev/null || true
}

# Help function
show_help() {
    echo "Product Service Database Migration Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  DB_HOST        Database host (default: localhost)"
    echo "  DB_PORT        Database port (default: 3306)"
    echo "  DB_USER        Database user (default: root)"
    echo "  DB_PASSWORD    Database password (default: password)"
    echo "  DB_NAME        Database name (default: edot_product)"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Run with default settings"
    echo "  DB_HOST=192.168.1.100 $0            # Run with custom host"
    echo "  DB_PASSWORD=mysecret $0              # Run with custom password"
}

# Handle command line arguments
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac