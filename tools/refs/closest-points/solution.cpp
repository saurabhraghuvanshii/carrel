#include <queue>
#include <utility>
#include <vector>
using namespace std;

// A max-heap of the k closest so far, by squared distance.
vector<int> closest(vector<vector<int>>& points, int k) {
    priority_queue<pair<int, int>> heap;
    for (int i = 0; i < (int)points.size(); i++) {
        heap.push({points[i][0] * points[i][0] + points[i][1] * points[i][1], i});
        if ((int)heap.size() > k) heap.pop();
    }
    vector<int> out;
    for (; !heap.empty(); heap.pop()) out.push_back(heap.top().second);
    return out;
}
