#include <array>
#include <string>
#include <vector>
using namespace std;

int countDistinct(string& word) {
    vector<array<int, 26>> next(1);
    next[0].fill(-1);
    for (size_t i = 0; i < word.size(); i++) {
        int cur = 0;
        for (size_t k = i; k < word.size(); k++) {
            int c = word[k] - 'a';
            if (next[cur][c] < 0) {
                next[cur][c] = next.size();
                next.emplace_back();
                next.back().fill(-1);
            }
            cur = next[cur][c];
        }
    }
    return next.size() - 1;
}
