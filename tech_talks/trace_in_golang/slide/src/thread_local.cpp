#include <iostream>
#include <thread>
using namespace std;

// START OMIT
thread_local int requestId=0;

void g(){
    std::cout<< "requestId" << requestId << endl;
}

void threadfunc(int id){
    requestId = id;
    g();
}

int main(){
    std::thread t1(threadfunc,1);
    std::thread t2(threadfunc,2);
    std::thread t3(threadfunc,3);

    t1.join();
    t2.join();
    t3.join();
}
// output:
// requestId1
// requestId2
// requestId3
// END OMIT