import 'package:shelf_blog_api/common/model/base_entity.dart';

class Favorite extends BaseEntity {
  final String userId;
  final String articleId;

  Favorite(
      {required String id,
      required this.userId,
      required this.articleId,
      required DateTime createdAt,
      required DateTime updatedAt,
      required DateTime? deletedAt})
      : super(
            id: id,
            createdAt: createdAt,
            updatedAt: updatedAt,
            deletedAt: deletedAt);
}
