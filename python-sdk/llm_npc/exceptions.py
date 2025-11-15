"""Custom exceptions for the LLM NPC SDK."""


class LLMNPCError(Exception):
    """Base exception for all LLM NPC SDK errors."""
    pass


class BackendConnectionError(LLMNPCError):
    """Raised when unable to connect to the backend."""
    pass


class BackendError(LLMNPCError):
    """Raised when the backend returns an error response."""
    
    def __init__(self, message: str, status_code: int = None):
        super().__init__(message)
        self.status_code = status_code


class NPCNotFoundError(LLMNPCError):
    """Raised when an NPC is not found."""
    pass


class SessionError(LLMNPCError):
    """Raised when there's an error with session management."""
    pass


class ToolRegistrationError(LLMNPCError):
    """Raised when tool registration fails."""
    pass

