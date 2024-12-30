# hardcover.app integration

[hardcover.app](https://hardcover.app/) reading tracking.


## Graphql API

On the [hardcover.app/account/api](https://hardcover.app/account/api) page you get an authentication token valid for 1-year. And a link to the graphql explorer for their endpoint [https://api.hardcover.app/v1/graphql](https://api.hardcover.app/v1/graphql)

Running a test query should get this result
```json
query Test {
  me {
    username
  }
}

---
{"data": {
 "me": [{
        "username": "timchurchard"
    }]
  }
}
```

---
So it seems they only support date not time for reading session. I guess I'll have to do more 24-hour reading sessions

```sql
query ListMyBooks {
  me {
    username
    user_books {
      user_book_status {
        id
        status
      }
      book {
        id
        title
        contributions {
          author {
            id
            name
          }
        }
      }
      user_book_reads {
        id
        started_at
        paused_at
        finished_at
      }
    }
  }
}

query FindBookByTitle {
  books(where: {title: {_eq: "Finders Keepers"}}, limit: 5) {
    id
    slug
    title
    description
    pages
    contributions {
      author_id
      author {
        id
        name
      }
    }
  }
}
```

---
The website uses the graphql. This is an example record reading

{"operationName":"UpsertDatesReadMutation","variables":{"userBookId":2648080,"datesRead":[{"id":1167309,"action":"update","started_at":"2024-05-14","finished_at":null,"reading_format_id":1,"edition_id":31296521},{"id":null,"action":"insert","started_at":"2024-05-14","finished_at":"2024-05-15","reading_format_id":1,"edition_id":null}]},"query":"fragment EditionInfoFragment on editions {\n  id\n  title\n  releaseDate: release_date\n  pages\n  audioSeconds: audio_seconds\n  readingFormatId: reading_format_id\n  usersCount: users_count\n  cachedImage: cached_image\n  language {\n    language\n    __typename\n  }\n  reading_format {\n    format\n    __typename\n  }\n  __typename\n}\n\nfragment UserBookReadFragment on user_book_reads {\n  id\n  userBookId: user_book_id\n  startedAt: started_at\n  finishedAt: finished_at\n  readingFormatId: reading_format_id\n  editionId: edition_id\n  edition {\n    ...EditionInfoFragment\n    __typename\n  }\n  __typename\n}\n\nfragment UserBookForButtonFragment on user_books {\n  id\n  bookId: book_id\n  userId: user_id\n  statusId: status_id\n  rating\n  progress\n  privacySettingId: privacy_setting_id\n  hasReview: has_review\n  datesRead: user_book_reads {\n    ...UserBookReadFragment\n    __typename\n  }\n  __typename\n}\n\nmutation UpsertDatesReadMutation($userBookId: Int!, $datesRead: [DatesReadInput]!) {\n  upsertResult: upsert_user_book_reads(\n    user_book_id: $userBookId\n    datesRead: $datesRead\n  ) {\n    error\n    userBook: user_book {\n      ...UserBookForButtonFragment\n      __typename\n    }\n    __typename\n  }\n}"}