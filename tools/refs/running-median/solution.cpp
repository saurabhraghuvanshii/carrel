#include <functional>
#include <queue>
#include <vector>
using namespace std;

class MedianFinder {
    priority_queue<int> low;
    priority_queue<int, vector<int>, greater<int>> high;

public:
    MedianFinder() {}

    void add(int x) {
        low.push(x);
        high.push(low.top());
        low.pop();
        if (high.size() > low.size()) {
            low.push(high.top());
            high.pop();
        }
    }

    double median() {
        if (low.size() > high.size()) return low.top();
        return (low.top() + (long long)high.top()) / 2.0;
    }
};
