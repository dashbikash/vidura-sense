db.getCollection('htmlpages').aggregate([
    {$match: {
      "lock_expiry":{"$ne":null} }
      },
])