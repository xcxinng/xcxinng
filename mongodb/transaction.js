
// 两个事务同时修改同个文档，产生冲突 得到报错：
// WriteConflict error: this operation conflicted with another operation. Please retry your operation or multi-document transaction.
//
// MongoDB没有使用阻塞锁(行锁/悲观锁，一直等待直至得到锁或者超时？)，
// 而是采用乐观锁，给有冲突的事务返回一个WriteConflictError, 然后继续尝试重试，直至修改成功或超时
//
var session1=db.getMongo().startSession();
var session2=db.getMongo().startSession();

var session1Collection = session1.getDatabase(db.getName()).transTest
var session2Collection = session2.getDatabase(db.getName()).transTest;

session1.startTransaction();
session2.startTransaction();

session1Collection.update({_id:1},{$set:{value:1}});
session2Collection.update({_id:1},{$set:{value:2}});

session1.commitTransaction();
session2.commitTransaction();
