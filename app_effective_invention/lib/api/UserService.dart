import 'dart:convert';

import 'package:app_effective_invention/types/user.dart';
import 'package:http/http.dart' as http;

class ApiService {
  final String baseUrl = "http://localhost:8080/users";

  Future<SingleUserResp> createUser(String name, String email, String password) async {
    final response = await http.post(
      Uri.parse("$baseUrl/new"),
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({
         "name": name,
          "email": "test3@example.com",
          "password": "password123"
          }),
    );
    return SingleUserResp.fromJson(json.decode(response.body));
  }

  Future<AuthResp> authUser(String email, String password) async {
    final response = await http.post(
      Uri.parse("$baseUrl/login"),
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({
        "email": email,
        "password": password
      })
    );
    return AuthResp.fromJson(json.decode(response.body));
  }

  Future<MultiUserResp> fetchUsers(String token) async {
    final response = await http.get(
      Uri.parse("$baseUrl/all"),
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer $token"
        },
      );
    if (response.statusCode == 200) {
      return MultiUserResp.fromJson(json.decode(response.body));
    } else {
      throw Exception('Error getting users!');
    }
  }

  Future<SingleUserResp> fetchUserById(String token, String id) async {
    final response = await http.get(
      Uri.parse("$baseUrl/$id"),
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer $token"
        },
      );
    if (response.statusCode == 200) {
      return SingleUserResp.fromJson(json.decode(response.body));
    } else {
      throw Exception('No user found with id: $id');
    }
  }

  Future<SingleUserResp> fetchUserByEmail(String token, String email) async {
    final response = await http.get(
      Uri.parse("$baseUrl/$email"),
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer $token"
        },
      );
    if (response.statusCode == 200) {
      return SingleUserResp.fromJson(json.decode(response.body));
    } else {
      throw Exception('No user found with id: $email');
    }
  }

  Future<MsgResp> updateUser(String token, String id, String name, String email) async {
    final response = await http.put(
      Uri.parse('$baseUrl/update'),
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer $token"
        },
      body: jsonEncode({
        'id': id,
        'name': name, 
        'email': email
        }),
    );
    return MsgResp.fromJson(json.decode(response.body));
  }

    Future<MsgResp> updateUserPassword(String token, String id, String password) async {
    final response = await http.put(
      Uri.parse('$baseUrl/update/password'),
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer $token"
        },
      body: jsonEncode({
        'id': id,
        'password': password, 
        }),
    );
    return MsgResp.fromJson(json.decode(response.body));
  }

  Future<MsgResp> deleteUser(String token, String id) async {
    final response = await http.delete(
      Uri.parse('$baseUrl/$id'),
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer $token"
        },
      );
    if (response.statusCode != 200) {
      throw Exception('Failed to delete post');
    }
    return MsgResp.fromJson(json.decode(response.body));
  }
}