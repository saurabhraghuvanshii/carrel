#include <queue>
#include <unordered_map>
#include <utility>
#include <vector>
using namespace std;

vector<int> mostFrequent(vector<int>& nums, int k) {
    unordered_map<int, int> count;
    for (int v : nums) count[v]++;
    // Keep the k most frequent in a small heap ordered by count.
    priority_queue<pair<int, int>, vector<pair<int, int>>, greater<>> heap;
    for (auto& [v, c] : count) {
        heap.push({c, v});
        if ((int)heap.size() > k) heap.pop();
    }
    vector<int> out;
    while (!heap.empty()) {
        out.push_back(heap.top().second);
        heap.pop();
    }
    return out;
}
