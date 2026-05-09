import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/note.dart';
import '../utils/constants.dart';

class ApiService {
  final http.Client _client = http.Client();

  Future<List<Note>> fetchNotes() async {
    final response = await _client.get(Uri.parse(ApiConstants.baseUrl));
    if (response.statusCode == 200) {
      final List<dynamic> data = json.decode(response.body);
      return data.map((json) => Note.fromJson(json)).toList();
    }
    throw Exception('Gagal mengambil catatan');
  }

  Future<Note> createNote(Note note) async {
    final response = await _client.post(
      Uri.parse(ApiConstants.baseUrl),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(note.toJson()),
    );
    if (response.statusCode == 201) {
      return Note.fromJson(json.decode(response.body));
    }
    throw Exception('Gagal membuat catatan');
  }

  Future<Note> updateNote(Note note) async {
    final response = await _client.put(
      Uri.parse('${ApiConstants.baseUrl}/${note.id}'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(note.toJson()),
    );
    if (response.statusCode == 200) {
      return Note.fromJson(json.decode(response.body));
    }
    throw Exception('Gagal memperbarui catatan');
  }

  Future<void> deleteNote(String id) async {
    final response = await _client.delete(
      Uri.parse('${ApiConstants.baseUrl}/$id'),
    );
    if (response.statusCode != 204) {
      throw Exception('Gagal menghapus catatan');
    }
  }

  void dispose() {
    _client.close();
  }
}
