#include <functional>
#include <queue>
#include <tuple>
#include <vector>
using namespace std;

vector<int> mergeAll(vector<vector<int>>& lists) {
    priority_queue<tuple<int, int, int>, vector<tuple<int, int, int>>, greater<>> heap;
    for (int i = 0; i < (int)lists.size(); i++)
        if (!lists[i].empty()) heap.push({lists[i][0], i, 0});
    vector<int> out;
    while (!heap.empty()) {
        auto [v, list, pos] = heap.top();
        heap.pop();
        out.push_back(v);
        if (pos + 1 < (int)lists[list].size()) heap.push({lists[list][pos + 1], list, pos + 1});
    }
    return out;
}
