#include <functional>
#include <queue>
#include <tuple>
#include <vector>
using namespace std;

vector<int> mergeOrder(vector<vector<int>>& logs, vector<int>& offsets, int limit) {
    typedef tuple<int, int, int> Entry;  // true time, server, position
    priority_queue<Entry, vector<Entry>, greater<Entry>> heap;
    for (int s = 0; s < (int)logs.size(); s++)
        if (!logs[s].empty()) heap.push({logs[s][0] + offsets[s], s, 0});
    vector<int> out;
    while ((int)out.size() < limit) {
        auto [at, s, pos] = heap.top();
        heap.pop();
        out.push_back(s);
        if (pos + 1 < (int)logs[s].size()) heap.push({logs[s][pos + 1] + offsets[s], s, pos + 1});
    }
    return out;
}
