#include <algorithm>
#include <queue>
#include <vector>
using namespace std;

int finishTime(int n, vector<int>& durations, vector<vector<int>>& deps) {
    vector<vector<int>> after(n);
    vector<int> waiting(n, 0), start(n, 0);
    for (auto& d : deps) {
        after[d[0]].push_back(d[1]);
        waiting[d[1]]++;
    }
    queue<int> ready;
    for (int j = 0; j < n; j++)
        if (waiting[j] == 0) ready.push(j);
    int done = 0, finish = 0;
    while (!ready.empty()) {
        int j = ready.front();
        ready.pop();
        done++;
        int end = start[j] + durations[j];
        finish = max(finish, end);
        for (int next : after[j]) {
            start[next] = max(start[next], end);
            if (--waiting[next] == 0) ready.push(next);
        }
    }
    return done == n ? finish : -1;
}
