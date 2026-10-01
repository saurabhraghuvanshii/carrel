#include <vector>
using namespace std;

class MinStack {
    vector<int> values, mins;

public:
    MinStack() {}

    void push(int x) {
        values.push_back(x);
        if (mins.empty() || x <= mins.back()) mins.push_back(x);
    }

    void pop() {
        if (values.back() == mins.back()) mins.pop_back();
        values.pop_back();
    }

    int top() { return values.back(); }

    int minimum() { return mins.back(); }
};
