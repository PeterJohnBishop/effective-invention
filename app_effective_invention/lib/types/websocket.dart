
class WSMessage {
  final String type;    // "broadcast", "private", etc.
  final String target;  // UserId for private messages
  final String content; // The actual message
  final String sender;  // Server populates this, but we can read it

  WSMessage({
    required this.type,
    this.target = "",
    required this.content,
    this.sender = "",
  });

  // Convert JSON from Server -> Dart Object
  factory WSMessage.fromJson(Map<String, dynamic> json) {
    return WSMessage(
      type: json['type'] ?? 'unknown',
      target: json['target'] ?? '',
      content: json['content'] ?? '',
      sender: json['sender'] ?? 'server',
    );
  }

  // Convert Dart Object -> JSON for Server
  Map<String, dynamic> toJson() {
    return {
      'type': type,
      'target': target,
      'content': content,
      // 'sender' is usually ignored by the server when reading
    };
  }
}