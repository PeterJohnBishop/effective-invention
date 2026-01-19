import 'dart:convert';

class UserFile {
    String user;
    String id;
    String filekey;
    int createdAt;

    UserFile({
        required this.user,
        required this.id,
        required this.filekey,
        required this.createdAt,
    });

    factory UserFile.fromRawJson(String str) => UserFile.fromJson(json.decode(str));

    String toRawJson() => json.encode(toJson());

    factory UserFile.fromJson(Map<String, dynamic> json) => UserFile(
        user: json["user"],
        id: json["id"],
        filekey: json["filekey"],
        createdAt: json["createdAt"],
    );

    Map<String, dynamic> toJson() => {
        "user": user,
        "id": id,
        "filekey": filekey,
        "createdAt": createdAt,
    };
}
