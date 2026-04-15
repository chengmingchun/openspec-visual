# OpenSpec Visualizer

A visual interface for OpenSpec - a specification-driven AI coding tool that helps manage AI-assisted development workflows.

## Overview

OpenSpec Visualizer provides a graphical interface to interact with OpenSpec, a tool that solves the "alignment problem" in AI coding by using specification-driven development (SDD). This application helps you visualize and manage OpenSpec workflows, specifications, and changes.

## Features

- **Visual Specification Management**: Browse and edit OpenSpec specifications through a user-friendly interface
- **Change Tracking**: Monitor active changes and their progress
- **AI Integration**: Connect with AI models for specification generation and code implementation
- **Real-time Updates**: See changes reflected immediately in the visual interface
- **Project Context Management**: Maintain project-specific configurations and rules

## Project Structure

```
openspec-visualizer/
├── frontend/              # Web interface (HTML/CSS/JS)
│   ├── index.html        # Main application page
│   ├── style.css         # Styling
│   └── main.js           # Frontend logic
├── openspec/             # OpenSpec configuration and examples
│   ├── specs/            # System specifications
│   ├── changes/          # Active changes
│   └── project.md        # Project documentation
├── spec.md               # Complete OpenSpec specification guide
├── main.go              # Application entry point
├── server.go            # HTTP server implementation
├── ai.go                # AI integration logic
├── fs.go                # File system operations
├── go.mod               # Go module dependencies
└── go.sum               # Go dependency checksums
```

## OpenSpec Specification

The `spec.md` file contains a comprehensive guide to OpenSpec based on the official user guide, covering:

- **Core Concepts**: Specs, Changes, Delta Specs
- **Workflow**: Proposal → Specs → Design → Tasks
- **Commands**: `/opsx:propose`, `/opsx:explore`, `/opsx:apply`, `/opsx:archive`
- **Best Practices**: Progressive strictness, one-change-one-responsibility
- **Configuration**: `config.yaml` setup and schema customization

## Getting Started

### Prerequisites
- Go 1.21 or later
- Node.js (for frontend development)
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/chengmingchun/openspec-visual.git
   cd openspec-visual
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o openspec-visualizer.exe
   ```

4. Run the application:
   ```bash
   ./openspec-visualizer.exe
   ```

5. Open your browser and navigate to `http://localhost:8080`

## Usage

1. **Initialize OpenSpec**: Use the interface to initialize OpenSpec in your project
2. **Create Specifications**: Define system behavior using the visual editor
3. **Manage Changes**: Track active changes and their progress
4. **Generate Code**: Use AI to generate code based on specifications
5. **Archive Changes**: Complete and archive finished changes

## OpenSpec Commands

The application supports the following OpenSpec commands:

### Core Commands
- `/opsx:propose` - Generate complete change artifacts
- `/opsx:explore` - Research without creating files
- `/opsx:apply` - Execute implementation tasks
- `/opsx:archive` - Archive completed changes

### Expanded Commands
- `/opsx:new` - Create change skeleton
- `/opsx:continue` - Generate next artifact
- `/opsx:ff` - Fast-forward generate all artifacts
- `/opsx:verify` - Validate implementation against specs
- `/opsx:sync` - Sync specs without archiving
- `/opsx:bulk-archive` - Archive multiple changes

## Configuration

Configure OpenSpec through `openspec/config.yaml`:

```yaml
schema: spec-driven

context: |
  Tech Stack: TypeScript, React 18, Node.js, PostgreSQL
  API Style: RESTful
  Test Framework: Vitest + React Testing Library
  Code Standards: ESLint

rules:
  proposal:
    - Must include rollback plan
    - Must specify affected module scope
  specs:
    - Use Given/When/Then format for test scenarios
```

## Development

### Backend (Go)
The backend is built with Go and provides:
- HTTP server for frontend communication
- File system operations for OpenSpec management
- AI model integration for specification processing

### Frontend (HTML/JS/CSS)
The frontend is a single-page application with:
- Real-time specification visualization
- Interactive change management
- Responsive design for different screen sizes

### Building from Source
```bash
# Build for production
go build -ldflags="-s -w" -o openspec-visualizer.exe

# Development mode with hot reload (requires air)
air
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- OpenSpec by Fission AI for the specification-driven development paradigm
- The AI coding community for inspiration and best practices

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request