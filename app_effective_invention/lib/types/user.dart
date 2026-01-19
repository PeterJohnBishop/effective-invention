import 'dart:convert';

class User {
    String id;
    String name;
    String email;
    int? createdAt;
    int? updatedAt;

    User({
        required this.id,
        required this.name,
        required this.email,
        required this.createdAt,
        required this.updatedAt,
    });

    factory User.fromRawJson(String str) => User.fromJson(json.decode(str));

    String toRawJson() => json.encode(toJson());

    factory User.fromJson(Map<String, dynamic> json) => User(
        id: json["id"],
        name: json["name"],
        email: json["email"],
        createdAt: json["createdAt"],
        updatedAt: json["updatedAt"],
    );

    Map<String, dynamic> toJson() => {
        "id": id,
        "name": name,
        "email": email,
        "createdAt": createdAt,
        "updatedAt": updatedAt,
    };
}

class SingleUserResp {
    String message;
    User user;

    SingleUserResp({
        required this.message,
        required this.user,
    });

    factory SingleUserResp.fromRawJson(String str) => SingleUserResp.fromJson(json.decode(str));

    String toRawJson() => json.encode(toJson());

    factory SingleUserResp.fromJson(Map<String, dynamic> json) => SingleUserResp(
        message: json["message"],
        user: User.fromJson(json["user"]),
    );

    Map<String, dynamic> toJson() => {
        "message": message,
        "user": user.toJson(),
    };
}

class AuthResp {
    String message;
    String refreshToken;
    String token;
    User user;

    AuthResp({
        required this.message,
        required this.refreshToken,
        required this.token,
        required this.user,
    });

    factory AuthResp.fromRawJson(String str) => AuthResp.fromJson(json.decode(str));

    String toRawJson() => json.encode(toJson());

    factory AuthResp.fromJson(Map<String, dynamic> json) => AuthResp(
        message: json["message"],
        refreshToken: json["refresh_token"],
        token: json["token"],
        user: User.fromJson(json["user"]),
    );

    Map<String, dynamic> toJson() => {
        "message": message,
        "refresh_token": refreshToken,
        "token": token,
        "user": user.toJson(),
    };
}

class MultiUserResp {
    String message;
    List<User> users;

    MultiUserResp({
        required this.message,
        required this.users,
    });

    factory MultiUserResp.fromRawJson(String str) => MultiUserResp.fromJson(json.decode(str));

    String toRawJson() => json.encode(toJson());

    factory MultiUserResp.fromJson(Map<String, dynamic> json) => MultiUserResp(
        message: json["message"],
        users: json["users"] != null 
        ? List<User>.from(json["users"].map((x) => User.fromJson(x)))
        : [],
    );

    Map<String, dynamic> toJson() => {
        "message": message,
        "users": List<dynamic>.from(users.map((x) => x.toJson())),
    };
}

class MsgResp {
    String message;

    MsgResp({
        required this.message,
    });

    factory MsgResp.fromRawJson(String str) => MsgResp.fromJson(json.decode(str));

    String toRawJson() => json.encode(toJson());

    factory MsgResp.fromJson(Map<String, dynamic> json) => MsgResp(
        message: json["message"],
    );

    Map<String, dynamic> toJson() => {
        "message": message,
    };
}

