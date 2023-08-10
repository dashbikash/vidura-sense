db.getCollection('htmlpages').aggregate([
    {$match: {
      "body":null }
      },
    {
        $sort: {
          "updated_on": 1
        }
    },
    {$project: {
        "_id":0,
        "url":1
    }},
    {$skip: 0*20},
    {
    $limit: 20
  }
])