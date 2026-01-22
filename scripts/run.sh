#!/usr/bin/env bash
set -e

if [[ -t 1 ]]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    CYAN='\033[0;36m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    CYAN=''
    NC=''
fi

print_header() {
    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${CYAN}→${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

detect_os() {
    case "$(uname -s)" in
        Linux*)     OS="Linux";;
        Darwin*)    OS="macOS";;
        CYGWIN*)    OS="Windows";;
        MINGW*)     OS="Windows";;
        MSYS*)      OS="Windows";;
        *)          OS="Unknown";;
    esac
}

# Check dependencies
check_dependencies() {
    print_info "Checking dependencies..."

    local missing_deps=0

    # Check Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed!"
        echo "         Install at: https://docs.docker.com/get-docker/"
        missing_deps=1
    else
        print_success "Docker found: $(docker --version | cut -d' ' -f3 | cut -d',' -f1)"
    fi

    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Docker Compose is not installed!"
        echo "         Install at: https://docs.docker.com/compose/install/"
        missing_deps=1
    else
        if command -v docker-compose &> /dev/null; then
            print_success "Docker Compose found: $(docker-compose --version | cut -d' ' -f4 | cut -d',' -f1)"
        else
            print_success "Docker Compose found (plugin): $(docker compose version | cut -d' ' -f4)"
        fi
    fi

    # Check Make (optional)
    if command -v make &> /dev/null; then
        print_success "Make found: $(make --version | head -n1 | cut -d' ' -f3)"
        HAS_MAKE=1
    else
        print_warning "Make not found (optional)"
        HAS_MAKE=0
    fi

    if [ $missing_deps -eq 1 ]; then
        echo ""
        print_error "Required dependencies missing. Please install and try again."
        exit 1
    fi

    echo ""
}

# Determine docker-compose command
get_docker_compose_cmd() {
    if command -v docker-compose &> /dev/null; then
        echo "docker-compose"
    else
        echo "docker compose"
    fi
}

# Clean up previous containers
cleanup() {
    print_info "Cleaning up previous containers..."

    local compose_cmd=$(get_docker_compose_cmd)

    $compose_cmd down --remove-orphans 2>/dev/null || true

    print_success "Cleanup completed"
    echo ""
}

# Start application
start_app() {
    print_info "Starting containers..."

    local compose_cmd=$(get_docker_compose_cmd)

    # Build and start
    $compose_cmd up --build -d

    print_success "Containers started"
    echo ""
}

# Wait for services to be ready
wait_for_services() {
    print_info "Waiting for services to initialize..."

    local max_retries=60
    local retry_count=0
    local sleep_interval=1

    # Wait for PostgreSQL
    print_info "Checking PostgreSQL..."
    while [ $retry_count -lt $max_retries ]; do
        if docker exec accounts-postgres pg_isready -U postgres -d accounts_db &> /dev/null; then
            print_success "PostgreSQL ready!"
            break
        fi
        retry_count=$((retry_count + 1))
        sleep $sleep_interval
    done

    if [ $retry_count -eq $max_retries ]; then
        print_error "Timeout waiting for PostgreSQL"
        print_info "Check the logs: docker logs accounts-postgres"
        exit 1
    fi

    echo ""

    # Wait for API
    print_info "Checking API..."
    retry_count=0

    while [ $retry_count -lt $max_retries ]; do
        if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
            print_success "API ready!"
            break
        fi
        retry_count=$((retry_count + 1))
        sleep $sleep_interval
    done

    if [ $retry_count -eq $max_retries ]; then
        print_error "Timeout waiting for API"
        print_info "Check the logs: docker logs api-accounts"
        exit 1
    fi

    echo ""
}

# Display final information
show_info() {
    print_header "✓ Application Started Successfully!"
    echo ""
    echo -e "${CYAN}Operating System:${NC} $OS"
    echo ""
    echo -e "${CYAN}📍 Available Endpoints:${NC}"
    echo "   • API:        http://localhost:8080"
    echo "   • Health:     http://localhost:8080/health"
    echo "   • Swagger:    http://localhost:8080/swagger/index.html"
    echo ""
    echo -e "${CYAN}📊 Useful Commands:${NC}"

    local compose_cmd=$(get_docker_compose_cmd)

    echo "   • View logs:       $compose_cmd logs -f"
    echo "   • View API logs:   docker logs -f api-accounts"
    echo "   • View DB logs:    docker logs -f accounts-postgres"
    echo "   • Stop:            $compose_cmd down"
    echo "   • Restart:         $compose_cmd restart"

    if [ $HAS_MAKE -eq 1 ]; then
        echo ""
        echo -e "${CYAN}🛠  Make Commands:${NC}"
        echo "   • View help:       make help"
        echo "   • Run tests:       make test"
        echo "   • View coverage:   make test-coverage-html"
        echo "   • Access DB:       make db-shell"
    fi

    echo ""
    echo -e "${CYAN}💡 Tip:${NC} To stop the application, press Ctrl+C or run: $compose_cmd down"
    echo ""
}

# Check if should follow logs
follow_logs() {
    local compose_cmd=$(get_docker_compose_cmd)

    echo ""
    read -p "Do you want to follow the logs? [y/N] " -n 1 -r
    echo ""

    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo ""
        print_info "Following logs (Ctrl+C to exit)..."
        echo ""
        $compose_cmd logs -f
    fi
}

# Main function
main() {
    clear

    print_header "🚀 Launcher"
    echo ""

    detect_os
    check_dependencies
    cleanup
    start_app
    wait_for_services
    show_info
    follow_logs
}

# Run
main