#include <functional>
#include <queue>
#include <vector>
using namespace std;

int joinCost(vector<int>& ropes) {
    priority_queue<int, vector<int>, greater<int>> heap(ropes.begin(), ropes.end());
    int cost = 0;
    while (heap.size() > 1) {
        int a = heap.top();
        heap.pop();
        int b = heap.top();
        heap.pop();
        cost += a + b;
        heap.push(a + b);
    }
    return cost;
}
