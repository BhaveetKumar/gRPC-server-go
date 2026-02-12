## gRPC-Server 

#### 
All API's are gRPC based API's.

####
**Objective:** Develop a gRPC-based API in Golang to manage blog posts for a hypothetical
blogging platform. The API should support CRUD operations for blog posts, and each post
should have the following attributes:
****
1. PostID (unique identifier)
2. Title
3. Content
4. Author
5. Publication Date
6. Tags (multiple tags per post)
****

####
API Supports the following operations:

**Create Post**
*Input*: Post details (Title, Content, Author, Publication Date, Tags)
*Output*: The Post (PostID, Title, Content, Author, Publication Date, Tags). Error message, if
creation fails.


**Read Post**
*Input*: PostID of the post to retrieve
*Output*: Post details (PostID, Title, Content, Author, Publication Date, Tags) or an error message
if the post is not found.


**Update Post**
**Input*:* PostID of the post to update and new details (Title, Content, Author, Tags)
**Output**: Post details (PostID, Title, Content, Author, Publication Date, Tags) or error message if
the update failed


**Delete Post**
*Input*: PostID of the post to delete
*Output*: Success/Failure message