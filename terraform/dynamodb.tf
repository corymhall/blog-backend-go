resource "aws_dynamodb_table" "posts" {
  name           = "Posts"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"
  range_key      = "posted_date"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "posted_date"
    type = "S"
  }

  attribute {
    name = "author"
    type = "S"
  }

  global_secondary_index {
    name            = "author-posted_date-index"
    hash_key        = "author"
    range_key       = "posted_date"
    read_capacity   = 5
    write_capacity  = 5
    projection_type = "ALL"
  }
}

resource "aws_dynamodb_table" "comments" {
  name           = "Comments"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "post_id"
  range_key      = "comment_date"

  attribute {
    name = "post_id"
    type = "S"
  }

  attribute {
    name = "comment_date"
    type = "S"
  }
}

resource "aws_dynamodb_table" "reply" {
  name           = "Reply"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"
  range_key      = "reply_date"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "reply_date"
    type = "S"
  }
}

resource "aws_dynamodb_table" "user" {
  name           = "User"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}
