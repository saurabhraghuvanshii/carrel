A news site logs the id of the page for every view. Return the `k` most viewed pages for the dashboard, the most viewed first. When two pages have the same number of views, the smaller id comes first. If the log has fewer than `k` different pages, return all of them.

Count the views per page in a hash map, then sort the pages by count (largest first) and by id (smallest first) for ties.

## Constraints

- 1 ≤ n ≤ 10⁴
- 1 ≤ k ≤ 5
- 1 ≤ page id ≤ 10⁶

## Follow-up

The log has billions of views and thousands of different pages, but `k` stays tiny. How would you avoid sorting every page, and what changes if the dashboard wants this for only the last hour?
