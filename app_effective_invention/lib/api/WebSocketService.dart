import 'dart:convert';
import 'package:app_effective_invention/types/websocket.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class WebSocketService {
  WebSocketChannel? _channel;

  // Connect to the server with the required User ID
  void connect(String userId) {
    // NOTE: Use '10.0.2.2' for Android Emulator, 'localhost' for iOS Simulator
    final String url = "ws://10.0.2.2:8080/ws?id=$userId";
    
    _channel = WebSocketChannel.connect(Uri.parse(url));
    print("Websocket connected for user: $userId");
  }

  // Send a message (Matches your Go ReadPump logic)
  void sendMessage(String type, String content, {String target = ""}) {
    if (_channel != null) {
      final msg = WSMessage(type: type, content: content, target: target);
      _channel!.sink.add(jsonEncode(msg.toJson()));
    }
  }

  // Expose the stream to the UI
  Stream<WSMessage> get messages {
    if (_channel == null) return const Stream.empty();
    
    return _channel!.stream.map((data) {
      final json = jsonDecode(data);
      return WSMessage.fromJson(json);
    });
  }

  void disconnect() {
    _channel?.sink.close();
  }
}